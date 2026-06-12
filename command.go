package buttplug

// CommandValue is either a float64 percent (0.0-1.0) or an int step value.
type CommandValue = any

// DeviceOutputCommand is a command to send to a device output.
type DeviceOutputCommand struct {
	OutputType OutputType
	Value      CommandValue
	Duration   *int
}
