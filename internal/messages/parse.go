package messages

import (
	"fmt"
	"strconv"
)

// ParseMessage parses a single message from protocol format.
func ParseMessage(data map[string]any) (Message, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf("expected single message type, got %d", len(data))
	}

	var msgType string
	var msgData map[string]any
	for k, v := range data {
		msgType = k
		var ok bool
		msgData, ok = v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid message data for type %s", msgType)
		}
		break
	}

	switch msgType {
	case "RequestServerInfo":
		return parseRequestServerInfo(msgData)
	case "ServerInfo":
		return parseServerInfo(msgData)
	case "Ok":
		return parseOk(msgData)
	case "Error":
		return parseError(msgData)
	case "Ping":
		return parsePing(msgData)
	case "Disconnect":
		return parseDisconnect(msgData)
	case "StartScanning":
		return parseStartScanning(msgData)
	case "StopScanning":
		return parseStopScanning(msgData)
	case "ScanningFinished":
		return parseScanningFinished(msgData)
	case "RequestDeviceList":
		return parseRequestDeviceList(msgData)
	case "DeviceList":
		return parseDeviceList(msgData)
	case "OutputCmd":
		return parseOutputCmd(msgData)
	case "InputCmd":
		return parseInputCmd(msgData)
	case "InputReading":
		return parseInputReading(msgData)
	case "StopCmd":
		return parseStopCmd(msgData)
	default:
		return nil, fmt.Errorf("unknown message type: %s", msgType)
	}
}

// ParseMessages parses an array of messages from protocol format.
func ParseMessages(data []map[string]any) ([]Message, error) {
	result := make([]Message, 0, len(data))
	for _, item := range data {
		msg, err := ParseMessage(item)
		if err != nil {
			return nil, err
		}
		result = append(result, msg)
	}
	return result, nil
}

func parseID(data map[string]any) (int, error) {
	id, err := asInt(data["Id"])
	if err != nil {
		return 0, fmt.Errorf("invalid Id: %w", err)
	}
	return id, nil
}

func parseRequestServerInfo(data map[string]any) (*RequestServerInfo, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	clientName, _ := data["ClientName"].(string)
	major, _ := asInt(data["ProtocolVersionMajor"])
	minor, _ := asInt(data["ProtocolVersionMinor"])
	if major == 0 {
		major = 4
	}
	return &RequestServerInfo{
		BaseMessage:          BaseMessage{ID: id},
		ClientName:           clientName,
		ProtocolVersionMajor: major,
		ProtocolVersionMinor: minor,
	}, nil
}

func parseServerInfo(data map[string]any) (*ServerInfo, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	maxPing, err := asInt(data["MaxPingTime"])
	if err != nil {
		return nil, err
	}
	major, err := asInt(data["ProtocolVersionMajor"])
	if err != nil {
		return nil, err
	}
	minor, err := asInt(data["ProtocolVersionMinor"])
	if err != nil {
		return nil, err
	}
	msg := &ServerInfo{
		BaseMessage:          BaseMessage{ID: id},
		MaxPingTime:          maxPing,
		ProtocolVersionMajor: major,
		ProtocolVersionMinor: minor,
	}
	if name, ok := data["ServerName"].(string); ok {
		msg.ServerName = &name
	}
	return msg, nil
}

func parseOk(data map[string]any) (*Ok, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	return &Ok{BaseMessage: BaseMessage{ID: id}}, nil
}

func parseError(data map[string]any) (*Error, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	errorMessage, _ := data["ErrorMessage"].(string)
	code, err := asInt(data["ErrorCode"])
	if err != nil {
		return nil, err
	}
	return &Error{
		BaseMessage:  BaseMessage{ID: id},
		ErrorMessage: errorMessage,
		ErrorCode:    code,
	}, nil
}

func parsePing(data map[string]any) (*Ping, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	return &Ping{BaseMessage: BaseMessage{ID: id}}, nil
}

func parseDisconnect(data map[string]any) (*Disconnect, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	return &Disconnect{BaseMessage: BaseMessage{ID: id}}, nil
}

func parseStartScanning(data map[string]any) (*StartScanning, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	return &StartScanning{BaseMessage: BaseMessage{ID: id}}, nil
}

func parseStopScanning(data map[string]any) (*StopScanning, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	return &StopScanning{BaseMessage: BaseMessage{ID: id}}, nil
}

func parseScanningFinished(data map[string]any) (*ScanningFinished, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	return &ScanningFinished{BaseMessage: BaseMessage{ID: id}}, nil
}

func parseRequestDeviceList(data map[string]any) (*RequestDeviceList, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	return &RequestDeviceList{BaseMessage: BaseMessage{ID: id}}, nil
}

func parseDeviceList(data map[string]any) (*DeviceList, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	devicesRaw, _ := data["Devices"].(map[string]any)
	devices := make(map[int]DeviceInfo)
	for key, value := range devicesRaw {
		idx, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("invalid device index %q: %w", key, err)
		}
		deviceMap, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid device data for index %s", key)
		}
		device, err := parseDeviceInfo(deviceMap)
		if err != nil {
			return nil, err
		}
		devices[idx] = device
	}
	return &DeviceList{
		BaseMessage: BaseMessage{ID: id},
		Devices:     devices,
	}, nil
}

func parseDeviceInfo(data map[string]any) (DeviceInfo, error) {
	deviceName, _ := data["DeviceName"].(string)
	deviceIndex, err := asInt(data["DeviceIndex"])
	if err != nil {
		return DeviceInfo{}, err
	}
	timingGap, _ := asInt(data["DeviceMessageTimingGap"])

	info := DeviceInfo{
		DeviceName:             deviceName,
		DeviceIndex:            deviceIndex,
		DeviceMessageTimingGap: timingGap,
		DeviceFeatures:         make(map[int]DeviceFeatureDefinition),
	}
	if displayName, ok := data["DeviceDisplayName"].(string); ok {
		info.DeviceDisplayName = &displayName
	}
	if featuresRaw, ok := data["DeviceFeatures"].(map[string]any); ok {
		for key, value := range featuresRaw {
			idx, err := strconv.Atoi(key)
			if err != nil {
				return DeviceInfo{}, fmt.Errorf("invalid feature index %q: %w", key, err)
			}
			featureMap, ok := value.(map[string]any)
			if !ok {
				return DeviceInfo{}, fmt.Errorf("invalid feature data for index %s", key)
			}
			feature, err := parseDeviceFeatureDefinition(featureMap)
			if err != nil {
				return DeviceInfo{}, err
			}
			info.DeviceFeatures[idx] = feature
		}
	}
	return info, nil
}

func parseDeviceFeatureDefinition(data map[string]any) (DeviceFeatureDefinition, error) {
	featureIndex, err := asInt(data["FeatureIndex"])
	if err != nil {
		return DeviceFeatureDefinition{}, err
	}
	defn := DeviceFeatureDefinition{
		FeatureIndex: featureIndex,
	}
	if desc, ok := data["FeatureDescription"].(string); ok {
		defn.FeatureDescription = &desc
	}
	if outputRaw, ok := data["Output"].(map[string]any); ok {
		defn.Output = make(map[string]FeatureOutputDefinition, len(outputRaw))
		for k, v := range outputRaw {
			outputMap, ok := v.(map[string]any)
			if !ok {
				return DeviceFeatureDefinition{}, fmt.Errorf("invalid output definition for %s", k)
			}
			outDef, err := parseFeatureOutputDefinition(outputMap)
			if err != nil {
				return DeviceFeatureDefinition{}, err
			}
			defn.Output[k] = outDef
		}
	}
	if inputRaw, ok := data["Input"].(map[string]any); ok {
		defn.Input = make(map[string]FeatureInputDefinition, len(inputRaw))
		for k, v := range inputRaw {
			inputMap, ok := v.(map[string]any)
			if !ok {
				return DeviceFeatureDefinition{}, fmt.Errorf("invalid input definition for %s", k)
			}
			inDef, err := parseFeatureInputDefinition(inputMap)
			if err != nil {
				return DeviceFeatureDefinition{}, err
			}
			defn.Input[k] = inDef
		}
	}
	return defn, nil
}

func parseFeatureOutputDefinition(data map[string]any) (FeatureOutputDefinition, error) {
	value, err := asIntPair(data["Value"])
	if err != nil {
		return FeatureOutputDefinition{}, err
	}
	defn := FeatureOutputDefinition{Value: value}
	if durationRaw, ok := data["Duration"]; ok {
		duration, err := asIntPair(durationRaw)
		if err != nil {
			return FeatureOutputDefinition{}, err
		}
		defn.Duration = &duration
	}
	return defn, nil
}

func parseFeatureInputDefinition(data map[string]any) (FeatureInputDefinition, error) {
	valuesRaw, ok := data["Value"].([]any)
	if !ok {
		return FeatureInputDefinition{}, fmt.Errorf("invalid input value list")
	}
	values := make([][2]int, 0, len(valuesRaw))
	for _, item := range valuesRaw {
		pair, err := asIntPair(item)
		if err != nil {
			return FeatureInputDefinition{}, err
		}
		values = append(values, pair)
	}
	commandsRaw, ok := data["Command"].([]any)
	if !ok {
		return FeatureInputDefinition{}, fmt.Errorf("invalid input command list")
	}
	commands := make([]string, 0, len(commandsRaw))
	for _, item := range commandsRaw {
		cmd, ok := item.(string)
		if !ok {
			return FeatureInputDefinition{}, fmt.Errorf("invalid input command")
		}
		commands = append(commands, cmd)
	}
	return FeatureInputDefinition{Value: values, Command: commands}, nil
}

func parseOutputCmd(data map[string]any) (*OutputCmd, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	deviceIndex, err := asInt(data["DeviceIndex"])
	if err != nil {
		return nil, err
	}
	featureIndex, err := asInt(data["FeatureIndex"])
	if err != nil {
		return nil, err
	}
	command, ok := data["Command"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid OutputCmd command")
	}
	return &OutputCmd{
		BaseMessage:  BaseMessage{ID: id},
		DeviceIndex:  deviceIndex,
		FeatureIndex: featureIndex,
		Command:      command,
	}, nil
}

func parseInputCmd(data map[string]any) (*InputCmd, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	deviceIndex, err := asInt(data["DeviceIndex"])
	if err != nil {
		return nil, err
	}
	featureIndex, err := asInt(data["FeatureIndex"])
	if err != nil {
		return nil, err
	}
	inputType, _ := data["Type"].(string)
	command, _ := data["Command"].(string)
	return &InputCmd{
		BaseMessage:  BaseMessage{ID: id},
		DeviceIndex:  deviceIndex,
		FeatureIndex: featureIndex,
		InputType:    inputType,
		Command:      command,
	}, nil
}

func parseInputReading(data map[string]any) (*InputReading, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	deviceIndex, err := asInt(data["DeviceIndex"])
	if err != nil {
		return nil, err
	}
	featureIndex, err := asInt(data["FeatureIndex"])
	if err != nil {
		return nil, err
	}
	readingRaw, ok := data["Reading"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid InputReading data")
	}
	reading := make(map[string]InputReadingValue, len(readingRaw))
	for k, v := range readingRaw {
		valueMap, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid reading value for %s", k)
		}
		value, err := asInt(valueMap["Value"])
		if err != nil {
			return nil, err
		}
		reading[k] = InputReadingValue{Value: value}
	}
	return &InputReading{
		BaseMessage:  BaseMessage{ID: id},
		DeviceIndex:  deviceIndex,
		FeatureIndex: featureIndex,
		Reading:      reading,
	}, nil
}

func parseStopCmd(data map[string]any) (*StopCmd, error) {
	id, err := parseID(data)
	if err != nil {
		return nil, err
	}
	msg := &StopCmd{
		BaseMessage: BaseMessage{ID: id},
		Inputs:      true,
		Outputs:     true,
	}
	if deviceIndex, err := asInt(data["DeviceIndex"]); err == nil {
		msg.DeviceIndex = &deviceIndex
	}
	if featureIndex, err := asInt(data["FeatureIndex"]); err == nil {
		msg.FeatureIndex = &featureIndex
	}
	if inputs, ok := data["Inputs"].(bool); ok {
		msg.Inputs = inputs
	}
	if outputs, ok := data["Outputs"].(bool); ok {
		msg.Outputs = outputs
	}
	return msg, nil
}

func asInt(value any) (int, error) {
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("expected int, got %T", value)
	}
}

func asIntPair(value any) ([2]int, error) {
	arr, ok := value.([]any)
	if !ok || len(arr) != 2 {
		return [2]int{}, fmt.Errorf("expected [2]int array")
	}
	a, err := asInt(arr[0])
	if err != nil {
		return [2]int{}, err
	}
	b, err := asInt(arr[1])
	if err != nil {
		return [2]int{}, err
	}
	return [2]int{a, b}, nil
}
