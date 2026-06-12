// Connection - Connect to a Buttplug server.
//
// This is the simplest possible Buttplug example. It connects to a
// Buttplug server (like Intiface Central) and shows connection status.
//
// Prerequisites:
//  1. Install Intiface Central: https://intiface.com/central/
//  2. Start Intiface Central and click "Start Server"
//  3. Run: go run ./examples/connection
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hirusha-adi/buttplug-go"
)

func main() {
	ctx := context.Background()
	client := buttplug.NewClient("Connection Example")

	defer func() {
		if client.Connected() {
			if err := client.Disconnect(ctx); err != nil {
				log.Printf("disconnect error: %v", err)
			}
			fmt.Println("Disconnected.")
		}
	}()

	fmt.Println("Connecting to server...")
	if err := client.Connect(ctx, "ws://127.0.0.1:12345"); err != nil {
		var bpErr *buttplug.ButtplugError
		if errors.As(err, &bpErr) {
			fmt.Printf("Failed to connect: %s\n", bpErr.Message)
		} else {
			fmt.Printf("Failed to connect: %v\n", err)
		}
		return
	}

	serverName := "(unknown)"
	if name := client.ServerName(); name != nil {
		serverName = *name
	}
	fmt.Printf("Connected to: %s\n", serverName)
	fmt.Println("Connection successful!")
}
