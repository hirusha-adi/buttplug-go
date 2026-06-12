package messages

// RequestServerInfo is the client identification message sent at connection start.
type RequestServerInfo struct {
	BaseMessage
	ClientName             string `json:"ClientName"`
	ProtocolVersionMajor   int    `json:"ProtocolVersionMajor"`
	ProtocolVersionMinor   int    `json:"ProtocolVersionMinor"`
}

func NewRequestServerInfo(id int, clientName string) *RequestServerInfo {
	return &RequestServerInfo{
		BaseMessage:          BaseMessage{ID: id},
		ClientName:           clientName,
		ProtocolVersionMajor: 4,
		ProtocolVersionMinor: 0,
	}
}

func (m *RequestServerInfo) MessageType() string { return "RequestServerInfo" }

func (m *RequestServerInfo) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{
		"Id":                   m.ID,
		"ClientName":           m.ClientName,
		"ProtocolVersionMajor": m.ProtocolVersionMajor,
		"ProtocolVersionMinor": m.ProtocolVersionMinor,
	})
}

// ServerInfo is the server identification response.
type ServerInfo struct {
	BaseMessage
	ServerName           *string `json:"ServerName,omitempty"`
	MaxPingTime          int     `json:"MaxPingTime"`
	ProtocolVersionMajor int     `json:"ProtocolVersionMajor"`
	ProtocolVersionMinor int     `json:"ProtocolVersionMinor"`
}

func (m *ServerInfo) MessageType() string { return "ServerInfo" }

func (m *ServerInfo) ToProtocol() []map[string]any {
	fields := map[string]any{
		"Id":                   m.ID,
		"MaxPingTime":          m.MaxPingTime,
		"ProtocolVersionMajor": m.ProtocolVersionMajor,
		"ProtocolVersionMinor": m.ProtocolVersionMinor,
	}
	if m.ServerName != nil {
		fields["ServerName"] = *m.ServerName
	}
	return wrapMessage(m.MessageType(), fields)
}

// Ok is a success response from the server.
type Ok struct {
	BaseMessage
}

func (m *Ok) MessageType() string { return "Ok" }

func (m *Ok) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{"Id": m.ID})
}

// Error is an error response from the server.
type Error struct {
	BaseMessage
	ErrorMessage string             `json:"ErrorMessage"`
	ErrorCode    int `json:"ErrorCode"`
}

func (m *Error) MessageType() string { return "Error" }

func (m *Error) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{
		"Id":           m.ID,
		"ErrorMessage": m.ErrorMessage,
		"ErrorCode":    m.ErrorCode,
	})
}

// Ping is a keepalive ping message.
type Ping struct {
	BaseMessage
}

func (m *Ping) MessageType() string { return "Ping" }

func (m *Ping) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{"Id": m.ID})
}

// Disconnect is a graceful disconnection request.
type Disconnect struct {
	BaseMessage
}

func (m *Disconnect) MessageType() string { return "Disconnect" }

func (m *Disconnect) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{"Id": m.ID})
}

// StartScanning starts scanning for devices.
type StartScanning struct {
	BaseMessage
}

func (m *StartScanning) MessageType() string { return "StartScanning" }

func (m *StartScanning) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{"Id": m.ID})
}

// StopScanning stops scanning for devices.
type StopScanning struct {
	BaseMessage
}

func (m *StopScanning) MessageType() string { return "StopScanning" }

func (m *StopScanning) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{"Id": m.ID})
}

// ScanningFinished notifies that scanning has completed.
type ScanningFinished struct {
	BaseMessage
}

func (m *ScanningFinished) MessageType() string { return "ScanningFinished" }

func (m *ScanningFinished) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{"Id": m.ID})
}

// RequestDeviceList requests the list of connected devices.
type RequestDeviceList struct {
	BaseMessage
}

func (m *RequestDeviceList) MessageType() string { return "RequestDeviceList" }

func (m *RequestDeviceList) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{"Id": m.ID})
}
