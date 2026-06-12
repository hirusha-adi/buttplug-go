package messages

import "strconv"

// FeatureOutputDefinition describes output capability for a feature.
type FeatureOutputDefinition struct {
	Value    [2]int  `json:"Value"`
	Duration *[2]int `json:"Duration,omitempty"`
}

// FeatureInputDefinition describes input capability for a feature.
type FeatureInputDefinition struct {
	Value   [][2]int `json:"Value"`
	Command []string `json:"Command"`
}

// DeviceFeatureDefinition defines a single device feature.
type DeviceFeatureDefinition struct {
	FeatureIndex       int                               `json:"FeatureIndex"`
	FeatureDescription *string                           `json:"FeatureDescription,omitempty"`
	Output             map[string]FeatureOutputDefinition `json:"Output,omitempty"`
	Input              map[string]FeatureInputDefinition  `json:"Input,omitempty"`
}

// DeviceInfo describes a single device.
type DeviceInfo struct {
	DeviceName              string                             `json:"DeviceName"`
	DeviceIndex             int                                `json:"DeviceIndex"`
	DeviceMessageTimingGap  int                                `json:"DeviceMessageTimingGap"`
	DeviceDisplayName       *string                            `json:"DeviceDisplayName,omitempty"`
	DeviceFeatures          map[int]DeviceFeatureDefinition    `json:"DeviceFeatures"`
}

// DeviceList is the list of connected devices.
type DeviceList struct {
	BaseMessage
	Devices map[int]DeviceInfo `json:"Devices"`
}

func (m *DeviceList) MessageType() string { return "DeviceList" }

func (m *DeviceList) ToProtocol() []map[string]any {
	devices := make(map[string]any, len(m.Devices))
	for idx, device := range m.Devices {
		devices[intKey(idx)] = deviceToMap(device)
	}
	return wrapMessage(m.MessageType(), map[string]any{
		"Id":      m.ID,
		"Devices": devices,
	})
}

func deviceToMap(device DeviceInfo) map[string]any {
	features := make(map[string]any, len(device.DeviceFeatures))
	for idx, feature := range device.DeviceFeatures {
		features[intKey(idx)] = featureToMap(feature)
	}
	result := map[string]any{
		"DeviceName":             device.DeviceName,
		"DeviceIndex":            device.DeviceIndex,
		"DeviceMessageTimingGap": device.DeviceMessageTimingGap,
		"DeviceFeatures":         features,
	}
	if device.DeviceDisplayName != nil {
		result["DeviceDisplayName"] = *device.DeviceDisplayName
	}
	return result
}

func featureToMap(feature DeviceFeatureDefinition) map[string]any {
	result := map[string]any{
		"FeatureIndex": feature.FeatureIndex,
	}
	if feature.FeatureDescription != nil {
		result["FeatureDescription"] = *feature.FeatureDescription
	}
	if feature.Output != nil {
		output := make(map[string]any, len(feature.Output))
		for k, v := range feature.Output {
			out := map[string]any{"Value": []int{v.Value[0], v.Value[1]}}
			if v.Duration != nil {
				out["Duration"] = []int{v.Duration[0], v.Duration[1]}
			}
			output[k] = out
		}
		result["Output"] = output
	}
	if feature.Input != nil {
		input := make(map[string]any, len(feature.Input))
		for k, v := range feature.Input {
			values := make([][]int, len(v.Value))
			for i, pair := range v.Value {
				values[i] = []int{pair[0], pair[1]}
			}
			input[k] = map[string]any{
				"Value":   values,
				"Command": v.Command,
			}
		}
		result["Input"] = input
	}
	return result
}

func intKey(n int) string {
	return strconv.Itoa(n)
}
