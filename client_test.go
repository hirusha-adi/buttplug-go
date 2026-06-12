package buttplug_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hirusha-adi/buttplug-go"
)

func TestClientName(t *testing.T) {
	client := buttplug.NewClient("My Test App")
	if client.Name() != "My Test App" {
		t.Fatalf("unexpected name: %s", client.Name())
	}
}

func TestClientInitialState(t *testing.T) {
	client := buttplug.NewClient("Test")
	if client.Connected() || client.ServerName() != nil || client.Scanning() || len(client.Devices()) != 0 {
		t.Fatal("unexpected initial state")
	}
}

func TestClientEventCallbacksNilByDefault(t *testing.T) {
	client := buttplug.NewClient("Test")
	if client.OnDeviceAdded != nil || client.OnDeviceRemoved != nil || client.OnScanningFinished != nil ||
		client.OnServerDisconnect != nil || client.OnError != nil {
		t.Fatal("expected nil callbacks")
	}
}

func TestClientSetEventCallbacks(t *testing.T) {
	client := buttplug.NewClient("Test")
	added := func(*buttplug.ButtplugDevice) {}
	removed := func(*buttplug.ButtplugDevice) {}
	scanning := func() {}
	disconnect := func() {}
	onError := func(error) {}

	client.OnDeviceAdded = added
	client.OnDeviceRemoved = removed
	client.OnScanningFinished = scanning
	client.OnServerDisconnect = disconnect
	client.OnError = onError

	if client.OnDeviceAdded == nil || client.OnDeviceRemoved == nil || client.OnScanningFinished == nil ||
		client.OnServerDisconnect == nil || client.OnError == nil {
		t.Fatal("expected callbacks to be set")
	}
}

func TestStartScanningNotConnected(t *testing.T) {
	client := buttplug.NewClient("Test")
	err := client.StartScanning(context.Background())
	var connectorErr *buttplug.ButtplugConnectorError
	if !errors.As(err, &connectorErr) {
		t.Fatalf("expected connector error, got %v", err)
	}
}

func TestStopScanningNotConnected(t *testing.T) {
	client := buttplug.NewClient("Test")
	err := client.StopScanning(context.Background())
	var connectorErr *buttplug.ButtplugConnectorError
	if !errors.As(err, &connectorErr) {
		t.Fatalf("expected connector error, got %v", err)
	}
}

func TestStopAllDevicesNotConnected(t *testing.T) {
	client := buttplug.NewClient("Test")
	err := client.StopAllDevices(context.Background())
	var connectorErr *buttplug.ButtplugConnectorError
	if !errors.As(err, &connectorErr) {
		t.Fatalf("expected connector error, got %v", err)
	}
}

func TestDisconnectWhenNotConnected(t *testing.T) {
	client := buttplug.NewClient("Test")
	if err := client.Disconnect(context.Background()); err != nil {
		t.Fatal(err)
	}
	if client.Connected() {
		t.Fatal("expected disconnected")
	}
}
