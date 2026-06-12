// Device Info - Inspect device capabilities.
//
// This example shows how to inspect device features and capabilities:
// - List all available features
// - Check output types (vibrate, rotate, position)
// - Check input types (battery, sensors)
// - Access individual motors on multi-motor devices
//
// Prerequisites:
//  1. Install Intiface Central: https://intiface.com/central/
//  2. Start Intiface Central and click "Start Server"
//  3. Have a supported device connected
//  4. Run: go run ./examples/device_info
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hirusha-adi/buttplug-go"
)

func main() {
	ctx := context.Background()
	client := buttplug.NewClient("Device Info Example")

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
		fmt.Printf("\n%s\n", strings.Repeat("=", 50))
		fmt.Printf("Device: %s\n", device.Name())
		fmt.Printf("Index: %d\n", device.Index())

		displayName := "(none)"
		if name := device.DisplayName(); name != nil {
			displayName = *name
		}
		fmt.Printf("Display Name: %s\n", displayName)
		fmt.Printf("Timing Gap: %dms\n", device.MessageTimingGap())
		fmt.Printf("%s\n", strings.Repeat("=", 50))

		features := device.Features()
		fmt.Printf("\nFeatures (%d):\n", len(features))
		for _, feature := range features {
			description := "(no description)"
			if desc := feature.Description(); desc != nil {
				description = *desc
			}
			fmt.Printf("\n  Feature %d: %s\n", feature.Index(), description)

			if outputs := feature.Outputs(); len(outputs) > 0 {
				fmt.Println("    Outputs:")
				for outputType := range outputs {
					outputEnum := buttplug.OutputType(outputType)
					valueRange := feature.StepRange(outputEnum)
					durationRange := feature.DurationRange(outputEnum)

					if valueRange != nil {
						fmt.Printf("      - %s: values [%d, %d]", outputType, valueRange[0], valueRange[1])
						if durationRange != nil {
							fmt.Printf(", duration [%d, %d]ms", durationRange[0], durationRange[1])
						}
						fmt.Println()
					}
				}
			}

			if inputs := feature.Inputs(); len(inputs) > 0 {
				fmt.Println("    Inputs:")
				for inputType, inputDef := range inputs {
					fmt.Printf("      - %s: commands %v\n", inputType, inputDef.Command)
				}
			}
		}

		vibrateFeatures := device.GetFeaturesWithOutput(buttplug.OutputTypeVibrate)
		if len(vibrateFeatures) > 1 {
			fmt.Printf("\nThis device has %d independent vibrators!\n", len(vibrateFeatures))
			fmt.Println("Use feature.RunOutput() to control them individually.")
		}
	}

	_ = client.Disconnect(ctx)
	fmt.Println("\nDone!")
}
