# buttplug-go

The (unofficial) Go implementation of the [Buttplug](https://buttplug.io) Intimate Hardware Control Protocol (v4) client — a 1:1 port of the official Python [`buttplug-py`](https://github.com/buttplugio/buttplug-py) library. This was translated with the assistance of generative AI tools, but has been thoroughly tested with [Intiface Central](https://intiface.com/#intiface-central)'s simulations.

<img width="2172" height="724" alt="image" src="https://github.com/user-attachments/assets/81bcd9cf-19d6-4a4b-9f56-1336d27e30ba" />

## Installation

```bash
go get github.com/hirusha-adi/buttplug-go
```

## Quick Start

1. **Install and start [Intiface Central](https://intiface.com/central/)** — This is the server that connects to your devices.

2. **Connect and control devices:**

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/hirusha-adi/buttplug-go"
)

func main() {
    ctx := context.Background()
    client := buttplug.NewClient("My App")

    if err := client.Connect(ctx, "ws://127.0.0.1:12345"); err != nil {
        log.Fatal(err)
    }

    if err := client.StartScanning(ctx); err != nil {
        log.Fatal(err)
    }
    time.Sleep(5 * time.Second)
    _ = client.StopScanning(ctx)

    for _, device := range client.Devices() {
        fmt.Printf("Found: %s\n", device.Name())

        if device.HasOutput(buttplug.OutputTypeVibrate) {
            _ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
                OutputType: buttplug.OutputTypeVibrate,
                Value:      0.5,
            })
            time.Sleep(2 * time.Second)
            _ = device.Stop(ctx, true, true)
        }
    }

    _ = client.Disconnect(ctx)
}
```

## Features

- **Simple API**: Unified `RunOutput()` method for all output types
- **Full Protocol Support**: Implements Buttplug protocol v4
- **Idiomatic Go**: Context-based I/O and strong typing throughout
- **Event Callbacks**: Get notified when devices connect/disconnect

## Device Control

```go
// Check device capabilities and send commands
if device.HasOutput(buttplug.OutputTypeVibrate) {
    _ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
        OutputType: buttplug.OutputTypeVibrate,
        Value:      0.75,
    })
}

if device.HasOutput(buttplug.OutputTypeRotate) {
    _ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
        OutputType: buttplug.OutputTypeRotate,
        Value:      0.5,
    })
}

if device.HasOutput(buttplug.OutputTypePositionWithDuration) {
    duration := 500
    _ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
        OutputType: buttplug.OutputTypePositionWithDuration,
        Value:      1.0,
        Duration:   &duration,
    })
}

// Read sensors
if device.HasInput(buttplug.InputTypeBattery) {
    battery, err := device.Battery(ctx)
    if err == nil {
        fmt.Printf("Battery: %.0f%%\n", battery*100)
    }
}

// Stop device
_ = device.Stop(ctx, true, true)
```

## Event Handling

```go
// Set up callbacks before connecting
client.OnDeviceAdded = func(d *buttplug.ButtplugDevice) {
    fmt.Printf("Connected: %s\n", d.Name())
}
client.OnDeviceRemoved = func(d *buttplug.ButtplugDevice) {
    fmt.Printf("Disconnected: %s\n", d.Name())
}
client.OnScanningFinished = func() {
    fmt.Println("Scan complete")
}
client.OnServerDisconnect = func() {
    fmt.Println("Server disconnected!")
}

// Callbacks can start goroutines for async-style handling
client.OnDeviceAdded = func(device *buttplug.ButtplugDevice) {
    if device.HasOutput(buttplug.OutputTypeVibrate) {
        go func() {
            _ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
                OutputType: buttplug.OutputTypeVibrate,
                Value:      0.25,
            })
        }()
    }
}
```

## Examples

See the [examples/](examples/) directory for more detailed examples:

- `application` — Complete application workflow
- `connection` — Connecting to a server
- `device_control` — Vibrate, rotate, and position commands
- `device_control_simulated_stroker` — Simulated stroker control
- `device_enumeration` — Discovering devices
- `device_info` — Inspecting device features
- `sensors` — Battery and signal strength
- `errors` — Error handling

Ported from the official Python examples in [`buttplug-py`](https://github.com/buttplugio/buttplug-py/tree/main/examples):

| Example | Run |
|---|---|
| [application](https://github.com/hirusha-adi/buttplug-go/tree/main/examples/application) | `go run ./examples/application` |
| [connection](https://github.com/hirusha-adi/buttplug-go/tree/main/examples/connection) | `go run ./examples/connection` |
| [device_control](https://github.com/hirusha-adi/buttplug-go/tree/main/examples/device_control) | `go run ./examples/device_control` |
| [device_control_simulated_stroker](https://github.com/hirusha-adi/buttplug-go/tree/main/examples/device_control_simulated_stroker) | `go run ./examples/device_control_simulated_stroker` |
| [device_enumeration](https://github.com/hirusha-adi/buttplug-go/tree/main/examples/device_enumeration) | `go run ./examples/device_enumeration` |
| [device_info](https://github.com/hirusha-adi/buttplug-go/tree/main/examples/device_info) | `go run ./examples/device_info` |
| [errors](https://github.com/hirusha-adi/buttplug-go/tree/main/examples/errors) | `go run ./examples/errors` |
| [sensors](https://github.com/hirusha-adi/buttplug-go/tree/main/examples/sensors) | `go run ./examples/sensors` |

All examples expect [Intiface Central](https://intiface.com/central/) running at `ws://127.0.0.1:12345`.

## Package layout

| Python ([`buttplug-py`](https://github.com/buttplugio/buttplug-py)) | Go ([`buttplug-go`](https://github.com/hirusha-adi/buttplug-go)) |
|---|---|
| `buttplug.client` | [`client.go`](client.go) |
| `buttplug.device` | [`device.go`](device.go) |
| `buttplug.feature` | [`feature.go`](feature.go) |
| `buttplug.command` | [`command.go`](command.go) |
| `buttplug.connector` | [`connector.go`](connector.go) |
| `buttplug.enums` | [`enums.go`](enums.go) |
| `buttplug.errors` | [`errors.go`](errors.go) |
| `buttplug._messages` | [`internal/messages/`](internal/messages/) |
| `buttplug._utils` | [`internal/utils/`](internal/utils/) |

## API notes

- Python `async`/`await` maps to Go `context.Context` on all I/O methods.
- Event callbacks (`OnDeviceAdded`, etc.) are synchronous functions; set them before calling `Connect`.
- `ButtplugClient` / `NewButtplugClient` aliases are provided for parity with the Python naming.

## Requirements

- Go 1.22+
- [Intiface Central](https://intiface.com/central/) or another Buttplug server

## Documentation

- [Buttplug Developer Guide](https://docs.buttplug.io)
- [Protocol Specification](https://docs.buttplug.io/docs/spec)
- [Intiface Central](https://intiface.com/central/) — recommended Buttplug server
- [Device compatibility list](https://iostindex.com)

## Support

- [Discord](https://discord.buttplug.io) — Community chat and support
- [GitHub Issues](https://github.com/hirusha-adi/buttplug-go/issues) — Bug reports and feature requests
- [Contributing](CONTRIBUTING.md)
- [Patreon](https://patreon.com/qdot) / [GitHub Sponsors](https://github.com/sponsors/qdot) — Support Buttplug development

## License

This project is licensed under the BSD 3-Clause License. See [LICENSE](LICENSE).

The Go client is a port of [`buttplug-py`](https://github.com/buttplugio/buttplug-py); the [Buttplug protocol](https://buttplug.io) is maintained by [Nonpolynomial Labs](https://nonpolynomial.com).
