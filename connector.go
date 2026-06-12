package buttplug

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hirusha-adi/buttplug-go/internal/messages"
	"github.com/hirusha-adi/buttplug-go/internal/utils"
)

const defaultSendTimeout = 30 * time.Second

// WebSocketConnector handles WebSocket connection and message correlation.
// It is not intended for direct use; Client handles this internally.
type WebSocketConnector struct {
	url            string
	conn           *websocket.Conn
	messageSorter  *utils.MessageSorter
	receiveDone    chan struct{}
	connected      bool
	mu             sync.Mutex

	onMessage    func(messages.Message)
	onDisconnect func()
}

// NewWebSocketConnector creates a connector for the given WebSocket URL.
func NewWebSocketConnector(url string) *WebSocketConnector {
	return &WebSocketConnector{
		url:           url,
		messageSorter: utils.NewMessageSorter(),
	}
}

// Connected reports whether the connector is currently connected.
func (c *WebSocketConnector) Connected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected && c.conn != nil
}

// SetMessageCallback sets the callback for unsolicited server messages (Id=0).
func (c *WebSocketConnector) SetMessageCallback(callback func(messages.Message)) {
	c.onMessage = callback
}

// SetDisconnectCallback sets the callback for disconnection events.
func (c *WebSocketConnector) SetDisconnectCallback(callback func()) {
	c.onDisconnect = callback
}

// Connect establishes the WebSocket connection.
func (c *WebSocketConnector) Connect(ctx context.Context) error {
	c.mu.Lock()
	if c.connected {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.url, nil)
	if err != nil {
		return &ButtplugConnectorError{ButtplugError{Message: fmt.Sprintf("Failed to connect to %s: %v", c.url, err)}}
	}

	c.mu.Lock()
	c.conn = conn
	c.connected = true
	c.receiveDone = make(chan struct{})
	c.mu.Unlock()

	go c.receiveLoop()
	return nil
}

// Disconnect closes the WebSocket connection.
func (c *WebSocketConnector) Disconnect() {
	c.mu.Lock()
	if !c.connected {
		c.mu.Unlock()
		return
	}
	c.connected = false
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()

	if conn != nil {
		_ = conn.Close()
	}

	<-c.receiveDone
	c.messageSorter.RejectAll(&ButtplugConnectorError{ButtplugError{Message: "Connection closed"}})
}

// Send sends a message and waits for the correlated response.
func (c *WebSocketConnector) Send(ctx context.Context, message messages.Message, timeout time.Duration) (messages.Message, error) {
	if timeout == 0 {
		timeout = defaultSendTimeout
	}

	c.mu.Lock()
	if c.conn == nil || !c.connected {
		c.mu.Unlock()
		return nil, &ButtplugConnectorError{ButtplugError{Message: "Not connected"}}
	}
	conn := c.conn
	c.mu.Unlock()

	msgID := c.messageSorter.GetNextID()
	message.SetID(msgID)

	protocolData := message.ToProtocol()
	jsonBytes, err := json.Marshal(protocolData)
	if err != nil {
		return nil, &ButtplugConnectorError{ButtplugError{Message: fmt.Sprintf("Failed to serialize message: %v", err)}}
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
		return nil, &ButtplugConnectorError{ButtplugError{Message: fmt.Sprintf("Failed to send message: %v", err)}}
	}

	return c.messageSorter.WaitForResponse(ctx, msgID, timeout)
}

// SendNoResponse sends a message without waiting for a response.
func (c *WebSocketConnector) SendNoResponse(message messages.Message) error {
	c.mu.Lock()
	if c.conn == nil || !c.connected {
		c.mu.Unlock()
		return &ButtplugConnectorError{ButtplugError{Message: "Not connected"}}
	}
	conn := c.conn
	c.mu.Unlock()

	if message.GetID() == 0 {
		message.SetID(c.messageSorter.GetNextID())
	}

	protocolData := message.ToProtocol()
	jsonBytes, err := json.Marshal(protocolData)
	if err != nil {
		return &ButtplugConnectorError{ButtplugError{Message: fmt.Sprintf("Failed to serialize message: %v", err)}}
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
		return &ButtplugConnectorError{ButtplugError{Message: fmt.Sprintf("Failed to send message: %v", err)}}
	}
	return nil
}

func (c *WebSocketConnector) receiveLoop() {
	defer close(c.receiveDone)

	for {
		c.mu.Lock()
		conn := c.conn
		connected := c.connected
		c.mu.Unlock()

		if !connected || conn == nil {
			break
		}

		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var data []map[string]any
		if err := json.Unmarshal(raw, &data); err != nil {
			continue
		}

		parsed, err := messages.ParseMessages(data)
		if err != nil {
			continue
		}

		for _, msg := range parsed {
			if msg.GetID() == 0 {
				if c.onMessage != nil {
					func() {
						defer func() { _ = recover() }()
						c.onMessage(msg)
					}()
				}
			} else {
				c.messageSorter.Resolve(msg.GetID(), msg)
			}
		}
	}

	c.mu.Lock()
	wasConnected := c.connected
	c.connected = false
	c.mu.Unlock()

	if wasConnected && c.onDisconnect != nil {
		func() {
			defer func() { _ = recover() }()
			c.onDisconnect()
		}()
	}
}
