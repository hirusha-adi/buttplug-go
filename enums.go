package buttplug

// OutputType identifies device output types for controlling actuators.
type OutputType string

const (
	OutputTypeVibrate              OutputType = "Vibrate"
	OutputTypeRotate                 OutputType = "Rotate"
	OutputTypeOscillate              OutputType = "Oscillate"
	OutputTypeConstrict              OutputType = "Constrict"
	OutputTypeSpray                  OutputType = "Spray"
	OutputTypeTemperature            OutputType = "Temperature"
	OutputTypeLED                    OutputType = "Led"
	OutputTypePosition               OutputType = "Position"
	OutputTypePositionWithDuration   OutputType = "HwPositionWithDuration"
)

// InputType identifies device input types for reading sensors.
type InputType string

const (
	InputTypeBattery  InputType = "Battery"
	InputTypeRSSI     InputType = "RSSI"
	InputTypePressure InputType = "Pressure"
	InputTypeButton   InputType = "Button"
)

// InputCommandType identifies commands for input sensors.
type InputCommandType string

const (
	InputCommandTypeRead        InputCommandType = "Read"
	InputCommandTypeSubscribe   InputCommandType = "Subscribe"
	InputCommandTypeUnsubscribe InputCommandType = "Unsubscribe"
)

// ErrorCode is a protocol error code.
type ErrorCode int

const (
	ErrorCodeUnknown ErrorCode = 0
	ErrorCodeInit    ErrorCode = 1
	ErrorCodePing    ErrorCode = 2
	ErrorCodeMsg     ErrorCode = 3
	ErrorCodeDevice  ErrorCode = 4
)
