// Device Enumeration - Scan for and list devices.
//
// This example shows how to scan for devices and handle device
// connection/disconnection events.
//
// Prerequisites:
//  1. Install Intiface Central: https://intiface.com/central/
//  2. Start Intiface Central and click "Start Server"
//  3. Have a supported device nearby and powered on
//  4. Run: go run ./examples/device_enumeration
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hirusha-adi/buttplug-go"
)

func onDeviceAdded(device *buttplug.ButtplugDevice) {
	fmt.Printf("Device connected: %s (index %d)\n", device.Name(), device.Index())
}

func onDeviceRemoved(device *buttplug.ButtplugDevice) {
	fmt.Printf("Device disconnected: %s\n", device.Name())
}

func main() {
	ctx := context.Background()
	client := buttplug.NewClient("Device Enumeration Example")

	client.OnDeviceAdded = onDeviceAdded
	client.OnDeviceRemoved = onDeviceRemoved

	fmt.Println("Connecting to server...")
	if err := client.Connect(ctx, "ws://127.0.0.1:12345"); err != nil {
		log.Fatal(err)
	}

	serverName := "(unknown)"
	if name := client.ServerName(); name != nil {
		serverName = *name
	}
	fmt.Printf("Connected to: %s\n", serverName)

	fmt.Println("\nScanning for devices (5 seconds)...")
	if err := client.StartScanning(ctx); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	if err := client.StopScanning(ctx); err != nil {
		log.Fatal(err)
	}

	devices := client.Devices()
	if len(devices) > 0 {
		fmt.Printf("\nFound %d device(s):\n", len(devices))
		for _, device := range devices {
			fmt.Printf("  - %s\n", device.Name())
			if displayName := device.DisplayName(); displayName != nil {
				fmt.Printf("    Display name: %s\n", *displayName)
			}
		}
	} else {
		fmt.Println("\nNo devices found.")
		fmt.Println("Make sure your device is on and in pairing mode.")
	}

	if err := client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nDone!")
}
