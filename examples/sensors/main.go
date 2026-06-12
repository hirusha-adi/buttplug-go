// Sensors - Read battery level and signal strength.
//
// This example shows how to read sensor data from devices:
// - Battery level (most Bluetooth devices)
// - RSSI (Bluetooth signal strength)
//
// Not all devices have sensors. The example checks what each device
// supports before trying to read.
//
// Prerequisites:
//  1. Install Intiface Central: https://intiface.com/central/
//  2. Start Intiface Central and click "Start Server"
//  3. Have a supported device connected
//  4. Run: go run ./examples/sensors
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
	client := buttplug.NewClient("Sensor Reading Example")

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
		fmt.Printf("\n%s:\n", device.Name())

		if device.HasInput(buttplug.InputTypeBattery) {
			battery, err := device.Battery(ctx)
			if err != nil {
				fmt.Printf("  Battery read failed: %v\n", err)
			} else {
				fmt.Printf("  Battery: %.0f%%\n", battery*100)
			}
		} else {
			fmt.Println("  No battery sensor")
		}

		if device.HasInput(buttplug.InputTypeRSSI) {
			rssi, err := device.RSSI(ctx)
			if err != nil {
				fmt.Printf("  RSSI read failed: %v\n", err)
			} else {
				quality := "Poor"
				switch {
				case rssi > -50:
					quality = "Excellent"
				case rssi > -70:
					quality = "Good"
				case rssi > -80:
					quality = "Fair"
				}
				fmt.Printf("  Signal: %d dBm (%s)\n", rssi, quality)
			}
		} else {
			fmt.Println("  No signal strength sensor")
		}
	}

	_ = client.Disconnect(ctx)
	fmt.Println("\nDone!")
}
