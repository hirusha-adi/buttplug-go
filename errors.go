package buttplug

import "fmt"

// ButtplugError is the base error for all Buttplug errors.
type ButtplugError struct {
	Message string
}

func (e *ButtplugError) Error() string {
	return e.Message
}

// ButtplugConnectorError indicates connection-related errors.
type ButtplugConnectorError struct {
	ButtplugError
}

// ButtplugHandshakeError indicates handshake failures.
type ButtplugHandshakeError struct {
	ButtplugError
}

// ButtplugPingError indicates ping timeout errors.
type ButtplugPingError struct {
	ButtplugError
}

// ButtplugDeviceError indicates device command failures.
type ButtplugDeviceError struct {
	ButtplugError
}

// ButtplugMessageError indicates message parsing or permission errors.
type ButtplugMessageError struct {
	ButtplugError
}

// ButtplugUnknownError indicates unknown server errors.
type ButtplugUnknownError struct {
	ButtplugError
}

// ErrorFromCode creates the appropriate error type from a protocol error code.
func ErrorFromCode(code ErrorCode, message string) error {
	switch code {
	case ErrorCodeInit:
		return &ButtplugHandshakeError{ButtplugError{Message: message}}
	case ErrorCodePing:
		return &ButtplugPingError{ButtplugError{Message: message}}
	case ErrorCodeMsg:
		return &ButtplugMessageError{ButtplugError{Message: message}}
	case ErrorCodeDevice:
		return &ButtplugDeviceError{ButtplugError{Message: message}}
	default:
		return &ButtplugUnknownError{ButtplugError{Message: message}}
	}
}

// NewConnectorError creates a connector error with a formatted message.
func NewConnectorError(format string, args ...any) error {
	return &ButtplugConnectorError{ButtplugError{Message: fmt.Sprintf(format, args...)}}
}
