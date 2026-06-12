// Simulated Stroker - Complex position pattern over 50 iterations.
//
// Drives a linear/stroker device through a repeating randomized sequence:
//   1. Always move to bottom
//   2. Randomly either move to middle, or move to top then bottom again
//
// Each move is followed by a 0.3 second pause. Runs 50 iterations, then stops.
//
// Prerequisites:
//  1. Install Intiface Central: https://intiface.com/central/
//  2. Start Intiface Central and click "Start Server"
//  3. Have a stroker/linear device connected
//  4. Run: go run ./examples/device_control_simulated_stroker
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hirusha-adi/buttplug-go"
)

const (
	iterations      = 50
	moveDelay       = 300 * time.Millisecond
	fullTravelMS    = 500
	middleTravelMS  = 250
	positionTop     = 1.0
	positionMiddle  = 0.5
	positionBottom  = 0.0
)

func main() {
	ctx := context.Background()
	client := buttplug.NewClient("Simulated Stroker Example")

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
	defer func() { _ = client.Disconnect(ctx) }()

	fmt.Println("Scanning for devices (5 seconds)...")
	if err := client.StartScanning(ctx); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	if err := client.StopScanning(ctx); err != nil {
		log.Fatal(err)
	}

	device := findStroker(client.Devices())
	if device == nil {
		log.Fatal("No stroker device found (needs HwPositionWithDuration output)")
	}

	fmt.Printf("\nRunning pattern on: %s\n", device.Name())
	fmt.Printf("%d iterations, %s pause between moves\n\n", iterations, moveDelay)

	move := func(label string, position float64, durationMS int) error {
		fmt.Printf("  -> %s (%.0f%%)\n", label, position*100)
		duration := durationMS
		if err := device.RunOutput(ctx, buttplug.DeviceOutputCommand{
			OutputType: buttplug.OutputTypePositionWithDuration,
			Value:      position,
			Duration:   &duration,
		}); err != nil {
			return err
		}
		time.Sleep(moveDelay)
		return nil
	}

	for i := 1; i <= iterations; i++ {
		fmt.Printf("Iteration %d/%d\n", i, iterations)

		if err := move("bottom", positionBottom, fullTravelMS); err != nil {
			log.Fatalf("move failed: %v", err)
		}

		if rand.Intn(2) == 0 {
			if err := move("middle", positionMiddle, middleTravelMS); err != nil {
				log.Fatalf("move failed: %v", err)
			}
		} else {
			if err := move("top", positionTop, fullTravelMS); err != nil {
				log.Fatalf("move failed: %v", err)
			}
			if err := move("bottom", positionBottom, fullTravelMS); err != nil {
				log.Fatalf("move failed: %v", err)
			}
		}

		fmt.Println()
	}

	fmt.Println("Stopping device...")
	if err := device.Stop(ctx, true, true); err != nil {
		log.Printf("stop failed: %v", err)
	}

	fmt.Println("Done!")
}

func findStroker(devices map[int]*buttplug.ButtplugDevice) *buttplug.ButtplugDevice {
	for _, device := range devices {
		if device.HasOutput(buttplug.OutputTypePositionWithDuration) {
			return device
		}
	}
	return nil
}
