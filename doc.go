// Package buttplug is a Go client library for the Buttplug protocol v4.
//
// Repository: https://github.com/hirusha-adi/buttplug-go
//
// Basic usage:
//
//	client := buttplug.NewClient("My App")
//	if err := client.Connect(ctx, "ws://127.0.0.1:12345"); err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Disconnect(ctx)
//
//	if err := client.StartScanning(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	time.Sleep(5 * time.Second)
//	_ = client.StopScanning(ctx)
//
//	for _, device := range client.Devices() {
//	    if device.HasOutput(buttplug.OutputTypeVibrate) {
//	        _ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
//	            OutputType: buttplug.OutputTypeVibrate,
//	            Value:      0.5,
//	        })
//	        time.Sleep(time.Second)
//	        _ = device.Stop(ctx, true, true)
//	    }
//	}
package buttplug

// Version is the library version.
const Version = "1.0.0"
