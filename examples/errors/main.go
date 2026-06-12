// Error Handling - Handle errors gracefully.
//
// This example shows how to handle various error conditions:
// - Connection failures
// - Device communication errors
// - Server disconnections
//
// Prerequisites:
//  1. Install Intiface Central: https://intiface.com/central/
//  2. Run: go run ./examples/errors
//     (server doesn't need to be running to see connection error handling)
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hirusha-adi/buttplug-go"
)

func main() {
	ctx := context.Background()
	client := buttplug.NewClient("Error Handling Example")

	client.OnServerDisconnect = func() {
		fmt.Println("Server disconnected unexpectedly!")
	}

	fmt.Println("Attempting to connect to server...")
	if err := client.Connect(ctx, "ws://127.0.0.1:12345"); err != nil {
		var connectorErr *buttplug.ButtplugConnectorError
		var handshakeErr *buttplug.ButtplugHandshakeError
		var bpErr *buttplug.ButtplugError

		switch {
		case errors.As(err, &connectorErr):
			fmt.Printf("Connection failed: %s\n", connectorErr.Message)
			fmt.Println("Is Intiface Central running?")
		case errors.As(err, &handshakeErr):
			fmt.Printf("Handshake failed: %s\n", handshakeErr.Message)
		case errors.As(err, &bpErr):
			fmt.Printf("Unexpected error: %s\n", bpErr.Message)
		default:
			fmt.Printf("Unexpected error: %v\n", err)
		}
		return
	}

	serverName := "(unknown)"
	if name := client.ServerName(); name != nil {
		serverName = *name
	}
	fmt.Printf("Connected to: %s\n", serverName)

	defer func() {
		if client.Connected() {
			_ = client.Disconnect(ctx)
			fmt.Println("\nDisconnected cleanly.")
		}
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Error during operation: %v\n", r)
			}
		}()

		fmt.Println("\nScanning for devices...")
		if err := client.StartScanning(ctx); err != nil {
			var pingErr *buttplug.ButtplugPingError
			if errors.As(err, &pingErr) {
				fmt.Println("Server ping timeout - connection lost")
				return
			}
			var bpErr *buttplug.ButtplugError
			if errors.As(err, &bpErr) {
				fmt.Printf("Error during operation: %s\n", bpErr.Message)
			}
			return
		}
		time.Sleep(3 * time.Second)
		_ = client.StopScanning(ctx)

		for _, device := range client.Devices() {
			fmt.Printf("\nControlling: %s\n", device.Name())
			if err := device.RunOutput(ctx, buttplug.DeviceOutputCommand{
				OutputType: buttplug.OutputTypeVibrate,
				Value:      0.5,
			}); err != nil {
				var deviceErr *buttplug.ButtplugDeviceError
				if errors.As(err, &deviceErr) {
					fmt.Printf("  Device error: %s\n", deviceErr.Message)
				} else {
					fmt.Printf("  Device error: %v\n", err)
				}
				continue
			}
			time.Sleep(time.Second)
			if err := device.Stop(ctx, true, true); err != nil {
				var deviceErr *buttplug.ButtplugDeviceError
				if errors.As(err, &deviceErr) {
					fmt.Printf("  Device error: %s\n", deviceErr.Message)
				}
				continue
			}
			fmt.Println("  Control successful!")
		}
	}()
}
