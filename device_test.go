package buttplug

import (
	"context"
	"errors"
	"testing"

	"github.com/hirusha-adi/buttplug-go/internal/messages"
)

type mockClient struct{}

func (mockClient) sendDeviceMessage(_ context.Context, _ messages.Message) (messages.Message, error) {
	return nil, nil
}

func vibrateFeature(client mockClient) *DeviceFeature {
	desc := "Main Motor"
	defn := messages.DeviceFeatureDefinition{
		FeatureIndex:       0,
		FeatureDescription: &desc,
		Output: map[string]messages.FeatureOutputDefinition{
			"Vibrate": {Value: [2]int{0, 20}},
		},
	}
	return NewDeviceFeature(client, 0, defn)
}

func batteryFeature(client mockClient) *DeviceFeature {
	desc := "Battery"
	defn := messages.DeviceFeatureDefinition{
		FeatureIndex:       1,
		FeatureDescription: &desc,
		Input: map[string]messages.FeatureInputDefinition{
			"Battery": {Value: [][2]int{{0, 100}}, Command: []string{"Read"}},
		},
	}
	return NewDeviceFeature(client, 0, defn)
}

func positionFeature(client mockClient) *DeviceFeature {
	desc := "Stroker"
	duration := [2]int{0, 1000}
	defn := messages.DeviceFeatureDefinition{
		FeatureIndex:       0,
		FeatureDescription: &desc,
		Output: map[string]messages.FeatureOutputDefinition{
			"HwPositionWithDuration": {Value: [2]int{0, 100}, Duration: &duration},
		},
	}
	return NewDeviceFeature(client, 0, defn)
}

func TestFeatureIndex(t *testing.T) {
	if vibrateFeature(mockClient{}).Index() != 0 {
		t.Fatal("unexpected index")
	}
}

func TestFeatureDescription(t *testing.T) {
	desc := vibrateFeature(mockClient{}).Description()
	if desc == nil || *desc != "Main Motor" {
		t.Fatalf("unexpected description: %+v", desc)
	}
}

func TestHasOutputTrue(t *testing.T) {
	feature := vibrateFeature(mockClient{})
	if !feature.HasOutput(OutputTypeVibrate) || !feature.HasOutputName("Vibrate") {
		t.Fatal("expected vibrate output")
	}
}

func TestHasOutputFalse(t *testing.T) {
	if vibrateFeature(mockClient{}).HasOutput(OutputTypeRotate) {
		t.Fatal("expected no rotate output")
	}
}

func TestHasInputTrue(t *testing.T) {
	feature := batteryFeature(mockClient{})
	if !feature.HasInput(InputTypeBattery) || !feature.HasInputName("Battery") {
		t.Fatal("expected battery input")
	}
}

func TestHasInputFalse(t *testing.T) {
	if batteryFeature(mockClient{}).HasInput(InputTypeRSSI) {
		t.Fatal("expected no rssi input")
	}
}

func TestSupportsInputCommand(t *testing.T) {
	feature := batteryFeature(mockClient{})
	if !feature.SupportsInputCommand(InputTypeBattery, InputCommandTypeRead) {
		t.Fatal("expected read support")
	}
	if feature.SupportsInputCommand(InputTypeBattery, InputCommandTypeSubscribe) {
		t.Fatal("expected no subscribe support")
	}
}

func TestStepRange(t *testing.T) {
	stepRange := vibrateFeature(mockClient{}).StepRange(OutputTypeVibrate)
	if stepRange == nil || *stepRange != [2]int{0, 20} {
		t.Fatalf("unexpected step range: %+v", stepRange)
	}
}

func TestStepRangeNotFound(t *testing.T) {
	if vibrateFeature(mockClient{}).StepRange(OutputTypeRotate) != nil {
		t.Fatal("expected nil step range")
	}
}

func TestStepCount(t *testing.T) {
	count := vibrateFeature(mockClient{}).StepCount(OutputTypeVibrate)
	if count == nil || *count != 20 {
		t.Fatalf("unexpected step count: %+v", count)
	}
}

func TestStepCountNotFound(t *testing.T) {
	if vibrateFeature(mockClient{}).StepCount(OutputTypeRotate) != nil {
		t.Fatal("expected nil step count")
	}
}

func TestDurationRange(t *testing.T) {
	duration := positionFeature(mockClient{}).DurationRange(OutputTypePositionWithDuration)
	if duration == nil || *duration != [2]int{0, 1000} {
		t.Fatalf("unexpected duration range: %+v", duration)
	}
}

func TestDurationRangeNotFound(t *testing.T) {
	if vibrateFeature(mockClient{}).DurationRange(OutputTypeVibrate) != nil {
		t.Fatal("expected nil duration range")
	}
}

func TestConvertToStepFloat(t *testing.T) {
	feature := vibrateFeature(mockClient{})
	cases := map[float64]int{0.5: 10, 0.0: 0, 1.0: 20, 0.05: 1}
	for value, want := range cases {
		got, err := feature.ConvertToStep(OutputTypeVibrate, value)
		if err != nil || got != want {
			t.Fatalf("convert %v: got %d want %d err %v", value, got, want, err)
		}
	}
}

func TestConvertToStepInt(t *testing.T) {
	feature := vibrateFeature(mockClient{})
	for _, value := range []int{10, 0, 20} {
		got, err := feature.ConvertToStep(OutputTypeVibrate, value)
		if err != nil || got != value {
			t.Fatalf("convert %d: got %d err %v", value, got, err)
		}
	}
}

func TestConvertToStepOutOfRangeFloat(t *testing.T) {
	_, err := vibrateFeature(mockClient{}).ConvertToStep(OutputTypeVibrate, 1.5)
	var deviceErr *ButtplugDeviceError
	if !errors.As(err, &deviceErr) {
		t.Fatalf("expected device error, got %v", err)
	}
}

func TestConvertToStepOutOfRangeInt(t *testing.T) {
	_, err := vibrateFeature(mockClient{}).ConvertToStep(OutputTypeVibrate, 25)
	var deviceErr *ButtplugDeviceError
	if !errors.As(err, &deviceErr) {
		t.Fatalf("expected device error, got %v", err)
	}
}

func multiFeatureDeviceInfo() messages.DeviceInfo {
	desc1, desc2, desc3, desc4 := "Vibrator 1", "Vibrator 2", "Rotator", "Battery"
	display := "My Device"
	return messages.DeviceInfo{
		DeviceName:             "Test Multi-Feature Device",
		DeviceIndex:            0,
		DeviceMessageTimingGap: 50,
		DeviceDisplayName:      &display,
		DeviceFeatures: map[int]messages.DeviceFeatureDefinition{
			0: {
				FeatureIndex: 0, FeatureDescription: &desc1,
				Output: map[string]messages.FeatureOutputDefinition{"Vibrate": {Value: [2]int{0, 20}}},
			},
			1: {
				FeatureIndex: 1, FeatureDescription: &desc2,
				Output: map[string]messages.FeatureOutputDefinition{"Vibrate": {Value: [2]int{0, 20}}},
			},
			2: {
				FeatureIndex: 2, FeatureDescription: &desc3,
				Output: map[string]messages.FeatureOutputDefinition{"Rotate": {Value: [2]int{0, 10}}},
			},
			3: {
				FeatureIndex: 3, FeatureDescription: &desc4,
				Input: map[string]messages.FeatureInputDefinition{
					"Battery": {Value: [][2]int{{0, 100}}, Command: []string{"Read"}},
				},
			},
		},
	}
}

func TestDeviceProperties(t *testing.T) {
	device := NewButtplugDevice(mockClient{}, multiFeatureDeviceInfo())
	if device.Index() != 0 || device.Name() != "Test Multi-Feature Device" {
		t.Fatal("unexpected device properties")
	}
	if device.DisplayName() == nil || *device.DisplayName() != "My Device" || device.MessageTimingGap() != 50 {
		t.Fatal("unexpected display/timing")
	}
	if len(device.Features()) != 4 {
		t.Fatalf("expected 4 features, got %d", len(device.Features()))
	}
}

func TestDeviceHasOutput(t *testing.T) {
	device := NewButtplugDevice(mockClient{}, multiFeatureDeviceInfo())
	if !device.HasOutput(OutputTypeVibrate) || !device.HasOutput(OutputTypeRotate) || device.HasOutput(OutputTypePosition) {
		t.Fatal("unexpected output capabilities")
	}
}

func TestDeviceHasInput(t *testing.T) {
	device := NewButtplugDevice(mockClient{}, multiFeatureDeviceInfo())
	if !device.HasInput(InputTypeBattery) || device.HasInput(InputTypeRSSI) {
		t.Fatal("unexpected input capabilities")
	}
}

func TestGetFeaturesWithOutput(t *testing.T) {
	device := NewButtplugDevice(mockClient{}, multiFeatureDeviceInfo())
	if len(device.GetFeaturesWithOutput(OutputTypeVibrate)) != 2 {
		t.Fatal("expected 2 vibrate features")
	}
	if len(device.GetFeaturesWithOutput(OutputTypeRotate)) != 1 {
		t.Fatal("expected 1 rotate feature")
	}
	if len(device.GetFeaturesWithOutput(OutputTypePosition)) != 0 {
		t.Fatal("expected 0 position features")
	}
}

func TestGetFeaturesWithInput(t *testing.T) {
	device := NewButtplugDevice(mockClient{}, multiFeatureDeviceInfo())
	if len(device.GetFeaturesWithInput(InputTypeBattery)) != 1 {
		t.Fatal("expected 1 battery feature")
	}
	if len(device.GetFeaturesWithInput(InputTypeRSSI)) != 0 {
		t.Fatal("expected 0 rssi features")
	}
}

func TestDeviceHasBattery(t *testing.T) {
	if !NewButtplugDevice(mockClient{}, multiFeatureDeviceInfo()).HasBattery() {
		t.Fatal("expected battery")
	}
}

func TestDeviceHasRSSI(t *testing.T) {
	if NewButtplugDevice(mockClient{}, multiFeatureDeviceInfo()).HasRSSI() {
		t.Fatal("expected no rssi")
	}
}

func TestDeviceFeatureStepCount(t *testing.T) {
	device := NewButtplugDevice(mockClient{}, multiFeatureDeviceInfo())
	vibrateCount := device.Features()[0].StepCount(OutputTypeVibrate)
	rotateCount := device.Features()[2].StepCount(OutputTypeRotate)
	if vibrateCount == nil || *vibrateCount != 20 || rotateCount == nil || *rotateCount != 10 {
		t.Fatal("unexpected step counts")
	}
}
