package messages

// Message is the base interface for all Buttplug protocol messages.
type Message interface {
	MessageType() string
	GetID() int
	SetID(id int)
	ToProtocol() []map[string]any
}

// BaseMessage holds the common Id field for all messages.
type BaseMessage struct {
	ID int `json:"Id"`
}

func (m *BaseMessage) GetID() int  { return m.ID }
func (m *BaseMessage) SetID(id int) { m.ID = id }

func wrapMessage(msgType string, fields map[string]any) []map[string]any {
	return []map[string]any{{msgType: fields}}
}
