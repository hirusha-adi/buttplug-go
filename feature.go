package buttplug

import (
	"context"
	"fmt"
	"math"

	"github.com/hirusha-adi/buttplug-go/internal/messages"
)

type deviceMessenger interface {
	sendDeviceMessage(ctx context.Context, msg messages.Message) (messages.Message, error)
}

// DeviceFeature represents a single feature of a device.
type DeviceFeature struct {
	client      deviceMessenger
	deviceIndex int
	definition  messages.DeviceFeatureDefinition
}

// NewDeviceFeature creates a feature from a protocol definition.
func NewDeviceFeature(client deviceMessenger, deviceIndex int, definition messages.DeviceFeatureDefinition) *DeviceFeature {
	return &DeviceFeature{
		client:      client,
		deviceIndex: deviceIndex,
		definition:  definition,
	}
}

// Index returns the feature index unique within the device.
func (f *DeviceFeature) Index() int {
	return f.definition.FeatureIndex
}

// Description returns the human-readable feature description.
func (f *DeviceFeature) Description() *string {
	return f.definition.FeatureDescription
}

// Outputs returns output types supported by this feature.
func (f *DeviceFeature) Outputs() map[string]messages.FeatureOutputDefinition {
	return f.definition.Output
}

// Inputs returns input types supported by this feature.
func (f *DeviceFeature) Inputs() map[string]messages.FeatureInputDefinition {
	return f.definition.Input
}

// HasOutput reports whether this feature supports a specific output type.
func (f *DeviceFeature) HasOutput(outputType OutputType) bool {
	if f.definition.Output == nil {
		return false
	}
	_, ok := f.definition.Output[string(outputType)]
	return ok
}

// HasOutputName reports whether this feature supports an output type by name.
func (f *DeviceFeature) HasOutputName(outputName string) bool {
	if f.definition.Output == nil {
		return false
	}
	_, ok := f.definition.Output[outputName]
	return ok
}

// HasInput reports whether this feature supports a specific input type.
func (f *DeviceFeature) HasInput(inputType InputType) bool {
	if f.definition.Input == nil {
		return false
	}
	_, ok := f.definition.Input[string(inputType)]
	return ok
}

// HasInputName reports whether this feature supports an input type by name.
func (f *DeviceFeature) HasInputName(inputName string) bool {
	if f.definition.Input == nil {
		return false
	}
	_, ok := f.definition.Input[inputName]
	return ok
}

// SupportsInputCommand reports whether this feature supports a specific input command.
func (f *DeviceFeature) SupportsInputCommand(inputType InputType, command InputCommandType) bool {
	if f.definition.Input == nil {
		return false
	}
	inputDef, ok := f.definition.Input[string(inputType)]
	if !ok {
		return false
	}
	for _, cmd := range inputDef.Command {
		if cmd == string(command) {
			return true
		}
	}
	return false
}

// StepRange returns the step value range for an output type.
func (f *DeviceFeature) StepRange(outputType OutputType) *[2]int {
	if f.definition.Output == nil {
		return nil
	}
	out, ok := f.definition.Output[string(outputType)]
	if !ok {
		return nil
	}
	return &out.Value
}

// StepCount returns the maximum step value for an output type.
func (f *DeviceFeature) StepCount(outputType OutputType) *int {
	stepRange := f.StepRange(outputType)
	if stepRange == nil {
		return nil
	}
	maxStep := stepRange[1]
	return &maxStep
}

// DurationRange returns the duration range for an output type.
func (f *DeviceFeature) DurationRange(outputType OutputType) *[2]int {
	if f.definition.Output == nil {
		return nil
	}
	out, ok := f.definition.Output[string(outputType)]
	if !ok {
		return nil
	}
	return out.Duration
}

// ConvertToStep converts a command value to a step value.
func (f *DeviceFeature) ConvertToStep(outputType OutputType, value CommandValue) (int, error) {
	stepRange := f.StepRange(outputType)
	if stepRange == nil {
		return 0, &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Feature does not support %s", outputType)}}
	}

	minStep, maxStep := stepRange[0], stepRange[1]
	var step int

	switch v := value.(type) {
	case float32:
		return convertFloatToStep(float64(v), minStep, maxStep)
	case float64:
		return convertFloatToStep(v, minStep, maxStep)
	case int:
		step = v
	case int64:
		step = int(v)
	default:
		return 0, &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Invalid command value type %T", value)}}
	}

	if step < minStep || step > maxStep {
		return 0, &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Step value %d out of range [%d, %d]", step, minStep, maxStep)}}
	}
	return step, nil
}

func convertFloatToStep(value float64, minStep, maxStep int) (int, error) {
	if value < -1.0 || value > 1.0 {
		return 0, &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Float value %v must be between -1.0 and 1.0", value)}}
	}

	var step int
	if value >= 0 {
		step = int(math.Ceil(value * float64(maxStep)))
	} else {
		step = int(math.Floor(value * float64(maxStep)))
	}

	if step < minStep || step > maxStep {
		return 0, &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Step value %d out of range [%d, %d]", step, minStep, maxStep)}}
	}
	return step, nil
}

// RunOutput sends an output command to this feature.
func (f *DeviceFeature) RunOutput(ctx context.Context, command DeviceOutputCommand) error {
	if command.OutputType == OutputTypePositionWithDuration {
		duration := 0
		if command.Duration != nil {
			duration = *command.Duration
		}
		return f.sendPositionWithDuration(ctx, command.Value, duration)
	}
	return f.sendOutput(ctx, command.OutputType, command.Value)
}

// Stop stops this feature's outputs.
func (f *DeviceFeature) Stop(ctx context.Context) error {
	featureIndex := f.Index()
	msg := &messages.StopCmd{
		BaseMessage:  messages.BaseMessage{ID: 0},
		DeviceIndex:  &f.deviceIndex,
		FeatureIndex: &featureIndex,
		Inputs:       false,
		Outputs:      true,
	}
	response, err := f.client.sendDeviceMessage(ctx, msg)
	if err != nil {
		return err
	}
	return f.checkResponse(response)
}

// Battery reads the battery level (0.0-1.0).
func (f *DeviceFeature) Battery(ctx context.Context) (float64, error) {
	reading, err := f.readInput(ctx, InputTypeBattery)
	if err != nil {
		return 0, err
	}
	return float64(reading) / 100.0, nil
}

// RSSI reads RSSI signal strength in dBm.
func (f *DeviceFeature) RSSI(ctx context.Context) (int, error) {
	return f.readInput(ctx, InputTypeRSSI)
}

func (f *DeviceFeature) sendOutput(ctx context.Context, outputType OutputType, value CommandValue) error {
	outputName := string(outputType)
	step, err := f.ConvertToStep(outputType, value)
	if err != nil {
		return err
	}

	msg := &messages.OutputCmd{
		BaseMessage:  messages.BaseMessage{ID: 0},
		DeviceIndex:  f.deviceIndex,
		FeatureIndex: f.Index(),
		Command: map[string]any{
			outputName: map[string]any{"Value": step},
		},
	}
	response, err := f.client.sendDeviceMessage(ctx, msg)
	if err != nil {
		return err
	}
	return f.checkResponse(response)
}

func (f *DeviceFeature) sendPositionWithDuration(ctx context.Context, value CommandValue, durationMS int) error {
	step, err := f.ConvertToStep(OutputTypePositionWithDuration, value)
	if err != nil {
		return err
	}

	if durationRange := f.DurationRange(OutputTypePositionWithDuration); durationRange != nil {
		minDur, maxDur := durationRange[0], durationRange[1]
		if durationMS < minDur {
			durationMS = minDur
		}
		if durationMS > maxDur {
			durationMS = maxDur
		}
	}

	msg := &messages.OutputCmd{
		BaseMessage:  messages.BaseMessage{ID: 0},
		DeviceIndex:  f.deviceIndex,
		FeatureIndex: f.Index(),
		Command: map[string]any{
			"HwPositionWithDuration": map[string]any{
				"Value":    step,
				"Duration": durationMS,
			},
		},
	}
	response, err := f.client.sendDeviceMessage(ctx, msg)
	if err != nil {
		return err
	}
	return f.checkResponse(response)
}

func (f *DeviceFeature) readInput(ctx context.Context, inputType InputType) (int, error) {
	inputName := string(inputType)
	if !f.HasInput(inputType) {
		return 0, &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Feature does not support %s input", inputName)}}
	}

	msg := &messages.InputCmd{
		BaseMessage:  messages.BaseMessage{ID: 0},
		DeviceIndex:  f.deviceIndex,
		FeatureIndex: f.Index(),
		InputType:    inputName,
		Command:      string(InputCommandTypeRead),
	}
	response, err := f.client.sendDeviceMessage(ctx, msg)
	if err != nil {
		return 0, err
	}

	if errMsg, ok := response.(*messages.Error); ok {
		return 0, ErrorFromCode(ErrorCode(errMsg.ErrorCode), errMsg.ErrorMessage)
	}

	reading, ok := response.(*messages.InputReading)
	if !ok {
		return 0, &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Unexpected response: %T", response)}}
	}

	value, ok := reading.Reading[inputName]
	if !ok {
		return 0, &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Invalid %s reading response", inputName)}}
	}
	return value.Value, nil
}

func (f *DeviceFeature) checkResponse(response messages.Message) error {
	if errMsg, ok := response.(*messages.Error); ok {
		return ErrorFromCode(ErrorCode(errMsg.ErrorCode), errMsg.ErrorMessage)
	}
	if _, ok := response.(*messages.Ok); !ok {
		return &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Unexpected response: %T", response)}}
	}
	return nil
}
