# buttplug-go

The (unofficial) Go implementation of the [Buttplug](https://buttplug.io) protocol v4 client — a 1:1 port of the official Python `buttplug-py` library.

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

| Python (`buttplug-py`) | Go (`buttplug-go`) |
|---|---|
| `buttplug.client` | `client.go` |
| `buttplug.device` | `device.go` |
| `buttplug.feature` | `feature.go` |
| `buttplug.command` | `command.go` |
| `buttplug.connector` | `connector.go` |
| `buttplug.enums` | `enums.go` |
| `buttplug.errors` | `errors.go` |
| `buttplug._messages` | `internal/messages/` |
| `buttplug._utils` | `internal/utils/` |

## API notes

- Python `async`/`await` maps to Go `context.Context` on all I/O methods.
- Event callbacks (`OnDeviceAdded`, etc.) are synchronous functions; set them before calling `Connect`.
- `ButtplugClient` / `NewButtplugClient` aliases are provided for parity with the Python naming.

## License

Same as the upstream Buttplug project.
