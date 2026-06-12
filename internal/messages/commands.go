package messages

// OutputCmd commands device output (vibrate, rotate, position, etc.).
type OutputCmd struct {
	BaseMessage
	DeviceIndex  int            `json:"DeviceIndex"`
	FeatureIndex int            `json:"FeatureIndex"`
	Command      map[string]any `json:"Command"`
}

func (m *OutputCmd) MessageType() string { return "OutputCmd" }

func (m *OutputCmd) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{
		"Id":           m.ID,
		"DeviceIndex":  m.DeviceIndex,
		"FeatureIndex": m.FeatureIndex,
		"Command":      m.Command,
	})
}

// InputCmd commands reading or subscribing to a device sensor.
type InputCmd struct {
	BaseMessage
	DeviceIndex  int    `json:"DeviceIndex"`
	FeatureIndex int    `json:"FeatureIndex"`
	InputType    string `json:"Type"`
	Command      string `json:"Command"`
}

func (m *InputCmd) MessageType() string { return "InputCmd" }

func (m *InputCmd) ToProtocol() []map[string]any {
	return wrapMessage(m.MessageType(), map[string]any{
		"Id":           m.ID,
		"DeviceIndex":  m.DeviceIndex,
		"FeatureIndex": m.FeatureIndex,
		"Type":         m.InputType,
		"Command":      m.Command,
	})
}

// InputReadingValue is a single input reading value.
type InputReadingValue struct {
	Value int `json:"Value"`
}

// InputReading is a sensor reading from a device.
type InputReading struct {
	BaseMessage
	DeviceIndex  int                          `json:"DeviceIndex"`
	FeatureIndex int                          `json:"FeatureIndex"`
	Reading      map[string]InputReadingValue `json:"Reading"`
}

func (m *InputReading) MessageType() string { return "InputReading" }

func (m *InputReading) ToProtocol() []map[string]any {
	reading := make(map[string]any, len(m.Reading))
	for k, v := range m.Reading {
		reading[k] = map[string]any{"Value": v.Value}
	}
	return wrapMessage(m.MessageType(), map[string]any{
		"Id":           m.ID,
		"DeviceIndex":  m.DeviceIndex,
		"FeatureIndex": m.FeatureIndex,
		"Reading":      reading,
	})
}

// StopCmd stops device outputs and/or unsubscribes from inputs.
type StopCmd struct {
	BaseMessage
	DeviceIndex  *int `json:"DeviceIndex,omitempty"`
	FeatureIndex *int `json:"FeatureIndex,omitempty"`
	Inputs       bool `json:"Inputs"`
	Outputs      bool `json:"Outputs"`
}

func NewStopCmd(id int) *StopCmd {
	return &StopCmd{
		BaseMessage: BaseMessage{ID: id},
		Inputs:      true,
		Outputs:     true,
	}
}

func (m *StopCmd) MessageType() string { return "StopCmd" }

func (m *StopCmd) ToProtocol() []map[string]any {
	fields := map[string]any{
		"Id":      m.ID,
		"Inputs":  m.Inputs,
		"Outputs": m.Outputs,
	}
	if m.DeviceIndex != nil {
		fields["DeviceIndex"] = *m.DeviceIndex
	}
	if m.FeatureIndex != nil {
		fields["FeatureIndex"] = *m.FeatureIndex
	}
	return wrapMessage(m.MessageType(), fields)
}
