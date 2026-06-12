# buttplug-go

The (unofficial) Go implementation of the [Buttplug](https://buttplug.io) protocol v4 client — a 1:1 port of the official Python [`buttplug-py`](https://github.com/buttplugio/buttplug-py) library.

**Repository:** [github.com/hirusha-adi/buttplug-go](https://github.com/hirusha-adi/buttplug-go)

<img width="2172" height="724" alt="image" src="https://github.com/user-attachments/assets/81bcd9cf-19d6-4a4b-9f56-1336d27e30ba" />

## Install

```bash
go get github.com/hirusha-adi/buttplug-go
```

## Usage

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/hirusha-adi/buttplug-go"
)

func main() {
    ctx := context.Background()
    client := buttplug.NewClient("My App")

    client.OnDeviceAdded = func(device *buttplug.ButtplugDevice) {
        log.Printf("Found: %s", device.Name())
    }

    if err := client.Connect(ctx, "ws://127.0.0.1:12345"); err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

    if err := client.StartScanning(ctx); err != nil {
        log.Fatal(err)
    }
    time.Sleep(5 * time.Second)
    _ = client.StopScanning(ctx)

    for _, device := range client.Devices() {
        if device.HasOutput(buttplug.OutputTypeVibrate) {
            _ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
                OutputType: buttplug.OutputTypeVibrate,
                Value:      0.5,
            })
            time.Sleep(time.Second)
            _ = device.Stop(ctx, true, true)
        }
    }
}
```

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

## Examples

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

## API notes

- Python `async`/`await` maps to Go `context.Context` on all I/O methods.
- Event callbacks (`OnDeviceAdded`, etc.) are synchronous functions; set them before calling `Connect`.
- `ButtplugClient` / `NewButtplugClient` aliases are provided for parity with the Python naming.

## Resources

- [Buttplug Developer Guide](https://docs.buttplug.io)
- [Protocol Specification](https://docs.buttplug.io/docs/spec)
- [Intiface Central](https://intiface.com/central/) — recommended Buttplug server
- [Device compatibility list](https://iostindex.com)
- [Issues & bug reports](https://github.com/hirusha-adi/buttplug-go/issues)
- [Contributing](CONTRIBUTING.md)

## License

This project is licensed under the BSD 3-Clause License. See [LICENSE](LICENSE).

The Go client is a port of [`buttplug-py`](https://github.com/buttplugio/buttplug-py); the [Buttplug protocol](https://buttplug.io) is maintained by [Nonpolynomial Labs](https://nonpolynomial.com).
