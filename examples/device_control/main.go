// Device Control - Vibrate, rotate, and position commands.
//
// This example shows how to control different types of devices:
// - Vibrators: Set vibration intensity
// - Rotators: Set rotation speed
// - Strokers: Move to position over time
//
// The example checks what each device supports before sending commands,
// so it will work with any device type.
//
// Prerequisites:
//  1. Install Intiface Central: https://intiface.com/central/
//  2. Start Intiface Central and click "Start Server"
//  3. Have a supported device connected
//  4. Run: go run ./examples/device_control
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
	client := buttplug.NewClient("Device Control Example")

	client.OnDeviceAdded = func(d *buttplug.ButtplugDevice) {
		fmt.Printf("Device connected: %s\n", d.Name())
	}
	client.OnDeviceRemoved = func(d *buttplug.ButtplugDevice) {
		fmt.Printf("Device disconnected: %s\n", d.Name())
	}

	fmt.Println("Connecting to server...")
	if err := client.Connect(ctx, "ws://127.0.0.1:12345"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scanning for devices (5 seconds)...")
	if err := client.StartScanning(ctx); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	if err := client.StopScanning(ctx); err != nil {
		log.Fatal(err)
	}

	devices := client.Devices()
	if len(devices) == 0 {
		fmt.Println("No devices found!")
		_ = client.Disconnect(ctx)
		return
	}

	for _, device := range devices {
		fmt.Printf("\nControlling: %s\n", device.Name())

		if device.HasOutput(buttplug.OutputTypeVibrate) {
			fmt.Println("  Starting vibration at 25%...")
			_ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
				OutputType: buttplug.OutputTypeVibrate,
				Value:      0.25,
			})
			time.Sleep(time.Second)

			fmt.Println("  Increasing to 50%...")
			_ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
				OutputType: buttplug.OutputTypeVibrate,
				Value:      0.5,
			})
			time.Sleep(time.Second)

			fmt.Println("  Full power (100%)...")
			_ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
				OutputType: buttplug.OutputTypeVibrate,
				Value:      1.0,
			})
			time.Sleep(time.Second)
		}

		if device.HasOutput(buttplug.OutputTypeRotate) {
			fmt.Println("  Rotating at 50%...")
			_ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
				OutputType: buttplug.OutputTypeRotate,
				Value:      0.5,
			})
			time.Sleep(2 * time.Second)
		}

		if device.HasOutput(buttplug.OutputTypePositionWithDuration) {
			duration500 := 500
			duration250 := 250

			fmt.Println("  Moving to top position...")
			_ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
				OutputType: buttplug.OutputTypePositionWithDuration,
				Value:      1.0,
				Duration:   &duration500,
			})
			time.Sleep(time.Second)

			fmt.Println("  Moving to bottom position...")
			_ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
				OutputType: buttplug.OutputTypePositionWithDuration,
				Value:      0.0,
				Duration:   &duration500,
			})
			time.Sleep(time.Second)

			fmt.Println("  Moving to middle...")
			_ = device.RunOutput(ctx, buttplug.DeviceOutputCommand{
				OutputType: buttplug.OutputTypePositionWithDuration,
				Value:      0.5,
				Duration:   &duration250,
			})
			time.Sleep(time.Second)
		}

		fmt.Println("  Stopping device...")
		_ = device.Stop(ctx, true, true)
	}

	fmt.Println("\nAll done!")
	_ = client.Disconnect(ctx)
}
