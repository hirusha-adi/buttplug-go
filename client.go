package buttplug

import (
	"context"
	"sync"
	"time"

	"github.com/hirusha-adi/buttplug-go/internal/messages"
)

// Client is the main client for connecting to Buttplug servers.
type Client struct {
	name    string
	mu      sync.Mutex
	connector *WebSocketConnector
	connected bool
	scanning  bool

	serverName  *string
	maxPingTime int

	pingCancel context.CancelFunc

	devices map[int]*ButtplugDevice

	OnDeviceAdded      func(*ButtplugDevice)
	OnDeviceRemoved    func(*ButtplugDevice)
	OnScanningFinished func()
	OnServerDisconnect func()
	OnError            func(error)
}

// NewClient creates a client with the given application name.
func NewClient(name string) *Client {
	return &Client{
		name:    name,
		devices: make(map[int]*ButtplugDevice),
	}
}

// Name returns the client application name.
func (c *Client) Name() string {
	return c.name
}

// Connected reports whether the client is connected to a server.
func (c *Client) Connected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

// ServerName returns the connected server name, or nil if not connected.
func (c *Client) ServerName() *string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.serverName
}

// Devices returns connected devices keyed by device index.
func (c *Client) Devices() map[int]*ButtplugDevice {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[int]*ButtplugDevice, len(c.devices))
	for k, v := range c.devices {
		copy[k] = v
	}
	return copy
}

// Scanning reports whether the client is currently scanning for devices.
func (c *Client) Scanning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.scanning
}

// Connect connects to a Buttplug server at the given WebSocket URL.
func (c *Client) Connect(ctx context.Context, url string) error {
	c.mu.Lock()
	if c.connected {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	connector := NewWebSocketConnector(url)
	connector.SetMessageCallback(c.handleServerMessage)
	connector.SetDisconnectCallback(c.handleDisconnect)

	if err := connector.Connect(ctx); err != nil {
		return err
	}

	request := messages.NewRequestServerInfo(0, c.name)
	response, err := connector.Send(ctx, request, 0)
	if err != nil {
		connector.Disconnect()
		return err
	}

	if errMsg, ok := response.(*messages.Error); ok {
		connector.Disconnect()
		return &ButtplugHandshakeError{ButtplugError{Message: errMsg.ErrorMessage}}
	}

	serverInfo, ok := response.(*messages.ServerInfo)
	if !ok {
		connector.Disconnect()
		return &ButtplugHandshakeError{ButtplugError{Message: "Unexpected response during handshake"}}
	}

	c.mu.Lock()
	c.connector = connector
	c.connected = true
	c.serverName = serverInfo.ServerName
	c.maxPingTime = serverInfo.MaxPingTime
	c.mu.Unlock()

	if c.maxPingTime > 0 {
		c.startPingTimer()
	}

	return c.requestDeviceList(ctx)
}

// Disconnect disconnects from the server and stops all devices.
func (c *Client) Disconnect(ctx context.Context) error {
	c.mu.Lock()
	if !c.connected || c.connector == nil {
		c.mu.Unlock()
		return nil
	}

	c.connected = false
	connector := c.connector
	c.connector = nil
	c.mu.Unlock()

	c.stopPingTimer()

	_ = c.StopAllDevices(ctx)

	connector.Disconnect()

	c.mu.Lock()
	c.devices = make(map[int]*ButtplugDevice)
	c.serverName = nil
	c.scanning = false
	c.mu.Unlock()

	return nil
}

// StartScanning starts scanning for devices.
func (c *Client) StartScanning(ctx context.Context) error {
	connector, err := c.requireConnector()
	if err != nil {
		return err
	}

	msg := &messages.StartScanning{BaseMessage: messages.BaseMessage{ID: 0}}
	response, err := connector.Send(ctx, msg, 0)
	if err != nil {
		return err
	}
	if errMsg, ok := response.(*messages.Error); ok {
		return ErrorFromCode(ErrorCode(errMsg.ErrorCode), errMsg.ErrorMessage)
	}

	c.mu.Lock()
	c.scanning = true
	c.mu.Unlock()
	return nil
}

// StopScanning stops scanning for devices.
func (c *Client) StopScanning(ctx context.Context) error {
	connector, err := c.requireConnector()
	if err != nil {
		return err
	}

	msg := &messages.StopScanning{BaseMessage: messages.BaseMessage{ID: 0}}
	response, err := connector.Send(ctx, msg, 0)
	if err != nil {
		return err
	}
	if errMsg, ok := response.(*messages.Error); ok {
		return ErrorFromCode(ErrorCode(errMsg.ErrorCode), errMsg.ErrorMessage)
	}

	c.mu.Lock()
	c.scanning = false
	c.mu.Unlock()
	return nil
}

// StopAllDevices stops all connected devices.
func (c *Client) StopAllDevices(ctx context.Context) error {
	connector, err := c.requireConnector()
	if err != nil {
		return err
	}

	msg := messages.NewStopCmd(0)
	response, err := connector.Send(ctx, msg, 0)
	if err != nil {
		return err
	}
	if errMsg, ok := response.(*messages.Error); ok {
		return ErrorFromCode(ErrorCode(errMsg.ErrorCode), errMsg.ErrorMessage)
	}
	return nil
}

func (c *Client) requestDeviceList(ctx context.Context) error {
	c.mu.Lock()
	connector := c.connector
	c.mu.Unlock()
	if connector == nil {
		return nil
	}

	msg := &messages.RequestDeviceList{BaseMessage: messages.BaseMessage{ID: 0}}
	response, err := connector.Send(ctx, msg, 0)
	if err != nil {
		return err
	}
	if deviceList, ok := response.(*messages.DeviceList); ok {
		return c.handleDeviceList(deviceList)
	}
	return nil
}

func (c *Client) handleServerMessage(msg messages.Message) {
	switch m := msg.(type) {
	case *messages.DeviceList:
		_ = c.handleDeviceList(m)
	case *messages.ScanningFinished:
		c.mu.Lock()
		c.scanning = false
		callback := c.OnScanningFinished
		c.mu.Unlock()
		if callback != nil {
			callback()
		}
	case *messages.Error:
		err := ErrorFromCode(ErrorCode(m.ErrorCode), m.ErrorMessage)
		c.mu.Lock()
		callback := c.OnError
		c.mu.Unlock()
		if callback != nil {
			callback(err)
		}
	}
}

func (c *Client) handleDeviceList(deviceList *messages.DeviceList) error {
	c.mu.Lock()
	current := make(map[int]struct{}, len(c.devices))
	for idx := range c.devices {
		current[idx] = struct{}{}
	}

	newIndices := make(map[int]struct{}, len(deviceList.Devices))
	for idx := range deviceList.Devices {
		newIndices[idx] = struct{}{}
	}

	var removed []*ButtplugDevice
	for idx := range current {
		if _, ok := newIndices[idx]; !ok {
			removed = append(removed, c.devices[idx])
			delete(c.devices, idx)
		}
	}

	onRemoved := c.OnDeviceRemoved
	for _, device := range removed {
		if onRemoved != nil {
			onRemoved(device)
		}
	}

	onAdded := c.OnDeviceAdded
	for idx := range newIndices {
		if _, exists := current[idx]; !exists {
			deviceInfo := deviceList.Devices[idx]
			device := NewButtplugDevice(c, deviceInfo)
			c.devices[idx] = device
			if onAdded != nil {
				onAdded(device)
			}
		}
	}
	c.mu.Unlock()
	return nil
}

func (c *Client) handleDisconnect() {
	c.mu.Lock()
	c.connected = false
	c.devices = make(map[int]*ButtplugDevice)
	c.scanning = false
	callback := c.OnServerDisconnect
	c.mu.Unlock()

	c.stopPingTimer()

	if callback != nil {
		callback()
	}
}

func (c *Client) startPingTimer() {
	c.stopPingTimer()

	ctx, cancel := context.WithCancel(context.Background())
	c.pingCancel = cancel

	go c.pingLoop(ctx)
}

func (c *Client) stopPingTimer() {
	if c.pingCancel != nil {
		c.pingCancel()
		c.pingCancel = nil
	}
}

func (c *Client) pingLoop(ctx context.Context) {
	c.mu.Lock()
	interval := time.Duration(c.maxPingTime) * time.Millisecond / 2
	c.mu.Unlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			connected := c.connected
			connector := c.connector
			onError := c.OnError
			c.mu.Unlock()

			if !connected || connector == nil {
				return
			}

			pingCtx, cancel := context.WithTimeout(ctx, defaultSendTimeout)
			msg := &messages.Ping{BaseMessage: messages.BaseMessage{ID: 0}}
			response, err := connector.Send(pingCtx, msg, defaultSendTimeout)
			cancel()
			if err != nil {
				continue
			}

			if errMsg, ok := response.(*messages.Error); ok {
				if onError != nil {
					onError(&ButtplugPingError{ButtplugError{Message: errMsg.ErrorMessage}})
				}
			}
		}
	}
}

func (c *Client) sendDeviceMessage(ctx context.Context, msg messages.Message) (messages.Message, error) {
	connector, err := c.requireConnector()
	if err != nil {
		return nil, err
	}
	return connector.Send(ctx, msg, 0)
}

func (c *Client) requireConnector() (*WebSocketConnector, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connector == nil || !c.connected {
		return nil, &ButtplugConnectorError{ButtplugError{Message: "Not connected"}}
	}
	return c.connector, nil
}

// ButtplugClient is an alias for Client matching the Python API name.
type ButtplugClient = Client

// NewButtplugClient is an alias for NewClient matching the Python API name.
func NewButtplugClient(name string) *Client {
	return NewClient(name)
}

// ButtplugDevice alias is defined in device.go as the struct name.
