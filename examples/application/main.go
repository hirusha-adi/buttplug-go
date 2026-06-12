// Buttplug Go - Complete Application Example
//
// This is a complete, working example that demonstrates the full workflow
// of a Buttplug application. If you're new to Buttplug, start here!
//
// Prerequisites:
//  1. Install Intiface Central: https://intiface.com/central
//  2. Start the server in Intiface Central (click "Start Server")
//  3. Run: go run ./examples/application
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hirusha-adi/buttplug-go"
)

func printDeviceCapabilities(device *buttplug.ButtplugDevice) {
	fmt.Printf("  %s\n", device.Name())

	var outputs []string
	if device.HasOutput(buttplug.OutputTypeVibrate) {
		outputs = append(outputs, "Vibrate")
	}
	if device.HasOutput(buttplug.OutputTypeRotate) {
		outputs = append(outputs, "Rotate")
	}
	if device.HasOutput(buttplug.OutputTypeOscillate) {
		outputs = append(outputs, "Oscillate")
	}
	if device.HasOutput(buttplug.OutputTypePosition) || device.HasOutput(buttplug.OutputTypePositionWithDuration) {
		outputs = append(outputs, "Position")
	}
	if device.HasOutput(buttplug.OutputTypeConstrict) {
		outputs = append(outputs, "Constrict")
	}
	if len(outputs) > 0 {
		fmt.Printf("    Outputs: %s\n", strings.Join(outputs, ", "))
	}

	var inputs []string
	if device.HasInput(buttplug.InputTypeBattery) {
		inputs = append(inputs, "Battery")
	}
	if device.HasInput(buttplug.InputTypeRSSI) {
		inputs = append(inputs, "RSSI")
	}
	if len(inputs) > 0 {
		fmt.Printf("    Inputs: %s\n", strings.Join(inputs, ", "))
	}

	fmt.Println()
}

func main() {
	ctx := context.Background()

	fmt.Println("===========================================")
	fmt.Println("  Buttplug Go Application Example")
	fmt.Println("===========================================")
	fmt.Println()

	client := buttplug.NewClient("My Buttplug Application")

	client.OnDeviceAdded = func(d *buttplug.ButtplugDevice) {
		fmt.Printf("[+] Device connected: %s\n", d.Name())
	}
	client.OnDeviceRemoved = func(d *buttplug.ButtplugDevice) {
		fmt.Printf("[-] Device disconnected: %s\n", d.Name())
	}
	client.OnServerDisconnect = func() {
		fmt.Println("[!] Server connection lost!")
	}

	fmt.Println("Connecting to Intiface Central...")
	if err := client.Connect(ctx, "ws://127.0.0.1:12345"); err != nil {
		var bpErr *buttplug.ButtplugError
		if errors.As(err, &bpErr) {
			fmt.Println("ERROR: Could not connect to Intiface Central!")
			fmt.Println("Make sure Intiface Central is running and the server is started.")
			fmt.Println("Default address: ws://127.0.0.1:12345")
			fmt.Printf("Error: %s\n", bpErr.Message)
		} else {
			fmt.Printf("ERROR: %v\n", err)
		}
		return
	}
	fmt.Println("Connected!")
	fmt.Println()

	fmt.Println("Scanning for devices...")
	fmt.Println("Turn on your Bluetooth/USB devices now.")
	fmt.Println()
	if err := client.StartScanning(ctx); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		_ = client.Disconnect(ctx)
		return
	}

	fmt.Print("Press Enter when your devices are connected...")
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	if err := client.StopScanning(ctx); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		_ = client.Disconnect(ctx)
		return
	}

	devices := client.Devices()
	if len(devices) == 0 {
		fmt.Println("No devices found. Make sure your device is:")
		fmt.Println("  - Turned on")
		fmt.Println("  - In pairing/discoverable mode")
		fmt.Println("  - Supported by Buttplug (check https://iostindex.com)")
		_ = client.Disconnect(ctx)
		return
	}

	fmt.Printf("\nFound %d device(s):\n\n", len(devices))

	deviceList := make([]*buttplug.ButtplugDevice, 0, len(devices))
	for _, device := range devices {
		deviceList = append(deviceList, device)
		printDeviceCapabilities(device)
	}

	fmt.Println("=== Interactive Control ===")
	fmt.Println("Commands:")
	fmt.Println("  v <0-100>  - Vibrate all devices at percentage")
	fmt.Println("  s          - Stop all devices")
	fmt.Println("  b          - Read battery levels")
	fmt.Println("  q          - Quit")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		userInput := strings.TrimSpace(strings.ToLower(line))
		if userInput == "" {
			continue
		}

		switch {
		case strings.HasPrefix(userInput, "v "):
			percentStr := strings.TrimSpace(userInput[2:])
			percent, err := strconv.Atoi(percentStr)
			if err != nil || percent < 0 || percent > 100 {
				fmt.Println("  Usage: v <0-100>")
				continue
			}
			intensity := float64(percent) / 100.0
			for _, device := range deviceList {
				if !device.HasOutput(buttplug.OutputTypeVibrate) {
					continue
				}
				if err := device.RunOutput(ctx, buttplug.DeviceOutputCommand{
					OutputType: buttplug.OutputTypeVibrate,
					Value:      intensity,
				}); err != nil {
					handleCommandError(err)
					continue
				}
				fmt.Printf("  %s: vibrating at %d%%\n", device.Name(), percent)
			}

		case userInput == "s":
			if err := client.StopAllDevices(ctx); err != nil {
				handleCommandError(err)
			} else {
				fmt.Println("  All devices stopped.")
			}

		case userInput == "b":
			for _, device := range deviceList {
				if device.HasInput(buttplug.InputTypeBattery) {
					battery, err := device.Battery(ctx)
					if err != nil {
						var deviceErr *buttplug.ButtplugDeviceError
						if errors.As(err, &deviceErr) {
							fmt.Printf("  %s: could not read battery - %s\n", device.Name(), deviceErr.Message)
						} else {
							fmt.Printf("  %s: could not read battery - %v\n", device.Name(), err)
						}
					} else {
						fmt.Printf("  %s: %.0f%% battery\n", device.Name(), battery*100)
					}
				} else {
					fmt.Printf("  %s: no battery sensor\n", device.Name())
				}
			}

		case userInput == "q":
			goto cleanup

		default:
			fmt.Println("  Unknown command. Use v, s, b, or q.")
		}
	}

cleanup:
	fmt.Println("\nStopping devices and disconnecting...")
	_ = client.StopAllDevices(ctx)
	_ = client.Disconnect(ctx)
	fmt.Println("Goodbye!")
}

func handleCommandError(err error) {
	var deviceErr *buttplug.ButtplugDeviceError
	var bpErr *buttplug.ButtplugError
	switch {
	case errors.As(err, &deviceErr):
		fmt.Printf("  Device error: %s\n", deviceErr.Message)
	case errors.As(err, &bpErr):
		fmt.Printf("  Error: %s\n", bpErr.Message)
	default:
		fmt.Printf("  Error: %v\n", err)
	}
}
