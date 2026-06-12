package buttplug

import (
	"context"
	"fmt"

	"github.com/hirusha-adi/buttplug-go/internal/messages"
)

// ButtplugDevice represents a connected Buttplug device.
type ButtplugDevice struct {
	client   deviceMessenger
	info     messages.DeviceInfo
	features map[int]*DeviceFeature
}

// NewButtplugDevice creates a device from protocol device info.
func NewButtplugDevice(client deviceMessenger, info messages.DeviceInfo) *ButtplugDevice {
	features := make(map[int]*DeviceFeature, len(info.DeviceFeatures))
	for idx, defn := range info.DeviceFeatures {
		features[idx] = NewDeviceFeature(client, info.DeviceIndex, defn)
	}
	return &ButtplugDevice{
		client:   client,
		info:     info,
		features: features,
	}
}

// Index returns the device index from the server.
func (d *ButtplugDevice) Index() int {
	return d.info.DeviceIndex
}

// Name returns the device name from server configuration.
func (d *ButtplugDevice) Name() string {
	return d.info.DeviceName
}

// DisplayName returns the user-provided display name, or nil if not set.
func (d *ButtplugDevice) DisplayName() *string {
	return d.info.DeviceDisplayName
}

// MessageTimingGap returns the minimum milliseconds between commands.
func (d *ButtplugDevice) MessageTimingGap() int {
	return d.info.DeviceMessageTimingGap
}

// Features returns device features keyed by feature index.
func (d *ButtplugDevice) Features() map[int]*DeviceFeature {
	return d.features
}

// HasOutput reports whether the device has any feature with the output type.
func (d *ButtplugDevice) HasOutput(outputType OutputType) bool {
	for _, feature := range d.features {
		if feature.HasOutput(outputType) {
			return true
		}
	}
	return false
}

// HasOutputName reports whether the device has any feature with the output name.
func (d *ButtplugDevice) HasOutputName(outputName string) bool {
	for _, feature := range d.features {
		if feature.HasOutputName(outputName) {
			return true
		}
	}
	return false
}

// HasInput reports whether the device has any feature with the input type.
func (d *ButtplugDevice) HasInput(inputType InputType) bool {
	for _, feature := range d.features {
		if feature.HasInput(inputType) {
			return true
		}
	}
	return false
}

// HasInputName reports whether the device has any feature with the input name.
func (d *ButtplugDevice) HasInputName(inputName string) bool {
	for _, feature := range d.features {
		if feature.HasInputName(inputName) {
			return true
		}
	}
	return false
}

// GetFeaturesWithOutput returns features that support a specific output type.
func (d *ButtplugDevice) GetFeaturesWithOutput(outputType OutputType) []*DeviceFeature {
	result := make([]*DeviceFeature, 0)
	for _, feature := range d.features {
		if feature.HasOutput(outputType) {
			result = append(result, feature)
		}
	}
	return result
}

// GetFeaturesWithInput returns features that support a specific input type.
func (d *ButtplugDevice) GetFeaturesWithInput(inputType InputType) []*DeviceFeature {
	result := make([]*DeviceFeature, 0)
	for _, feature := range d.features {
		if feature.HasInput(inputType) {
			result = append(result, feature)
		}
	}
	return result
}

// RunOutput sends an output command to all features matching the output type.
func (d *ButtplugDevice) RunOutput(ctx context.Context, command DeviceOutputCommand) error {
	features := d.GetFeaturesWithOutput(command.OutputType)
	if len(features) == 0 {
		return &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Device has no %s features", command.OutputType)}}
	}

	for _, feature := range features {
		if err := feature.RunOutput(ctx, command); err != nil {
			return err
		}
	}
	return nil
}

// Stop stops device outputs and/or unsubscribes from inputs.
func (d *ButtplugDevice) Stop(ctx context.Context, inputs, outputs bool) error {
	deviceIndex := d.Index()
	msg := &messages.StopCmd{
		BaseMessage: messages.BaseMessage{ID: 0},
		DeviceIndex: &deviceIndex,
		Inputs:      inputs,
		Outputs:     outputs,
	}
	response, err := d.client.sendDeviceMessage(ctx, msg)
	if err != nil {
		return err
	}
	return d.checkResponse(response)
}

// HasBattery reports whether the device has a battery sensor.
func (d *ButtplugDevice) HasBattery() bool {
	return d.HasInput(InputTypeBattery)
}

// Battery reads battery level from the first battery sensor (0.0-1.0).
func (d *ButtplugDevice) Battery(ctx context.Context) (float64, error) {
	features := d.GetFeaturesWithInput(InputTypeBattery)
	if len(features) == 0 {
		return 0, &ButtplugDeviceError{ButtplugError{Message: "Device has no battery sensor"}}
	}
	return features[0].Battery(ctx)
}

// HasRSSI reports whether the device has an RSSI sensor.
func (d *ButtplugDevice) HasRSSI() bool {
	return d.HasInput(InputTypeRSSI)
}

// RSSI reads RSSI from the first RSSI sensor.
func (d *ButtplugDevice) RSSI(ctx context.Context) (int, error) {
	features := d.GetFeaturesWithInput(InputTypeRSSI)
	if len(features) == 0 {
		return 0, &ButtplugDeviceError{ButtplugError{Message: "Device has no RSSI sensor"}}
	}
	return features[0].RSSI(ctx)
}

func (d *ButtplugDevice) checkResponse(response messages.Message) error {
	if errMsg, ok := response.(*messages.Error); ok {
		return ErrorFromCode(ErrorCode(errMsg.ErrorCode), errMsg.ErrorMessage)
	}
	if _, ok := response.(*messages.Ok); !ok {
		return &ButtplugDeviceError{ButtplugError{Message: fmt.Sprintf("Unexpected response: %T", response)}}
	}
	return nil
}
