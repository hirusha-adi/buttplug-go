package buttplug_test

import (
	"encoding/json"
	"testing"

	"github.com/hirusha-adi/buttplug-go"
	"github.com/hirusha-adi/buttplug-go/internal/messages"
)

func sampleDeviceListData() map[string]any {
	return map[string]any{
		"Id": 1,
		"Devices": map[string]any{
			"0": map[string]any{
				"DeviceName":             "Test Vibrator",
				"DeviceIndex":            0,
				"DeviceMessageTimingGap": 50,
				"DeviceDisplayName":      "My Vibrator",
				"DeviceFeatures": map[string]any{
					"0": map[string]any{
						"FeatureIndex":       0,
						"FeatureDescription": "Clitoral Stimulator",
						"Output":             map[string]any{"Vibrate": map[string]any{"Value": []any{0, 20}}},
					},
					"1": map[string]any{
						"FeatureIndex":       1,
						"FeatureDescription": "G-Spot Motor",
						"Output":             map[string]any{"Vibrate": map[string]any{"Value": []any{0, 20}}},
					},
					"2": map[string]any{
						"FeatureIndex":       2,
						"FeatureDescription": "Battery",
						"Input": map[string]any{
							"Battery": map[string]any{
								"Value":   []any{[]any{0, 100}},
								"Command": []any{"Read"},
							},
						},
					},
				},
			},
			"1": map[string]any{
				"DeviceName":             "Test Stroker",
				"DeviceIndex":            1,
				"DeviceMessageTimingGap": 100,
				"DeviceFeatures": map[string]any{
					"0": map[string]any{
						"FeatureIndex":       0,
						"FeatureDescription": "Stroker",
						"Output": map[string]any{
							"HwPositionWithDuration": map[string]any{
								"Value":    []any{0, 100},
								"Duration": []any{0, 1000},
							},
						},
					},
				},
			},
		},
	}
}

func TestRequestServerInfoSerialize(t *testing.T) {
	msg := messages.NewRequestServerInfo(1, "Test Client")
	result := msg.ToProtocol()

	expected := []map[string]any{{
		"RequestServerInfo": map[string]any{
			"Id":                   1,
			"ClientName":           "Test Client",
			"ProtocolVersionMajor": 4,
			"ProtocolVersionMinor": 0,
		},
	}}

	assertEqualMaps(t, result, expected)
}

func TestRequestServerInfoDeserialize(t *testing.T) {
	data := map[string]any{
		"RequestServerInfo": map[string]any{
			"Id":                   1,
			"ClientName":           "My App",
			"ProtocolVersionMajor": 4,
			"ProtocolVersionMinor": 0,
		},
	}
	msg, err := messages.ParseMessage(data)
	if err != nil {
		t.Fatal(err)
	}

	parsed, ok := msg.(*messages.RequestServerInfo)
	if !ok {
		t.Fatalf("expected RequestServerInfo, got %T", msg)
	}
	if parsed.GetID() != 1 || parsed.ClientName != "My App" || parsed.ProtocolVersionMajor != 4 {
		t.Fatalf("unexpected parsed message: %+v", parsed)
	}
}

func TestServerInfoDeserialize(t *testing.T) {
	data := map[string]any{
		"ServerInfo": map[string]any{
			"Id":                   1,
			"ServerName":           "Intiface Central",
			"MaxPingTime":          1000,
			"ProtocolVersionMajor": 4,
			"ProtocolVersionMinor": 0,
		},
	}
	msg, err := messages.ParseMessage(data)
	if err != nil {
		t.Fatal(err)
	}

	parsed := msg.(*messages.ServerInfo)
	if parsed.ServerName == nil || *parsed.ServerName != "Intiface Central" || parsed.MaxPingTime != 1000 {
		t.Fatalf("unexpected server info: %+v", parsed)
	}
}

func TestOkSerialize(t *testing.T) {
	msg := &messages.Ok{BaseMessage: messages.BaseMessage{ID: 5}}
	result := msg.ToProtocol()
	expected := []map[string]any{{"Ok": map[string]any{"Id": 5}}}
	assertEqualMaps(t, result, expected)
}

func TestErrorDeserialize(t *testing.T) {
	data := map[string]any{
		"Error": map[string]any{
			"Id":           0,
			"ErrorMessage": "Ping timeout",
			"ErrorCode":    2,
		},
	}
	msg, err := messages.ParseMessage(data)
	if err != nil {
		t.Fatal(err)
	}

	parsed := msg.(*messages.Error)
	if parsed.ErrorMessage != "Ping timeout" || parsed.ErrorCode != int(buttplug.ErrorCodePing) {
		t.Fatalf("unexpected error message: %+v", parsed)
	}
}

func TestPingRoundtrip(t *testing.T) {
	msg := &messages.Ping{BaseMessage: messages.BaseMessage{ID: 10}}
	protocol := msg.ToProtocol()
	parsed, err := messages.ParseMessage(protocol[0])
	if err != nil {
		t.Fatal(err)
	}
	if parsed.GetID() != 10 {
		t.Fatalf("expected id 10, got %d", parsed.GetID())
	}
}

func TestStartScanningSerialize(t *testing.T) {
	msg := &messages.StartScanning{BaseMessage: messages.BaseMessage{ID: 2}}
	result := msg.ToProtocol()
	expected := []map[string]any{{"StartScanning": map[string]any{"Id": 2}}}
	assertEqualMaps(t, result, expected)
}

func TestStopScanningSerialize(t *testing.T) {
	msg := &messages.StopScanning{BaseMessage: messages.BaseMessage{ID: 3}}
	result := msg.ToProtocol()
	expected := []map[string]any{{"StopScanning": map[string]any{"Id": 3}}}
	assertEqualMaps(t, result, expected)
}

func TestScanningFinishedDeserialize(t *testing.T) {
	data := map[string]any{"ScanningFinished": map[string]any{"Id": 0}}
	msg, err := messages.ParseMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	if msg.GetID() != 0 {
		t.Fatalf("expected id 0, got %d", msg.GetID())
	}
}

func TestRequestDeviceListSerialize(t *testing.T) {
	msg := &messages.RequestDeviceList{BaseMessage: messages.BaseMessage{ID: 4}}
	result := msg.ToProtocol()
	expected := []map[string]any{{"RequestDeviceList": map[string]any{"Id": 4}}}
	assertEqualMaps(t, result, expected)
}

func TestDeviceListDeserialize(t *testing.T) {
	data := map[string]any{"DeviceList": sampleDeviceListData()}
	msg, err := messages.ParseMessage(data)
	if err != nil {
		t.Fatal(err)
	}

	parsed := msg.(*messages.DeviceList)
	if parsed.GetID() != 1 || len(parsed.Devices) != 2 {
		t.Fatalf("unexpected device list: %+v", parsed)
	}

	device0 := parsed.Devices[0]
	if device0.DeviceName != "Test Vibrator" || device0.DeviceIndex != 0 || device0.DeviceMessageTimingGap != 50 {
		t.Fatalf("unexpected device0: %+v", device0)
	}
	if device0.DeviceDisplayName == nil || *device0.DeviceDisplayName != "My Vibrator" {
		t.Fatalf("unexpected display name: %+v", device0.DeviceDisplayName)
	}
	if len(device0.DeviceFeatures) != 3 {
		t.Fatalf("expected 3 features, got %d", len(device0.DeviceFeatures))
	}

	feature0 := device0.DeviceFeatures[0]
	if feature0.FeatureIndex != 0 || feature0.FeatureDescription == nil || *feature0.FeatureDescription != "Clitoral Stimulator" {
		t.Fatalf("unexpected feature0: %+v", feature0)
	}
	if feature0.Output["Vibrate"].Value != [2]int{0, 20} {
		t.Fatalf("unexpected vibrate range: %+v", feature0.Output["Vibrate"].Value)
	}

	feature2 := device0.DeviceFeatures[2]
	if feature2.Input["Battery"].Value[0] != [2]int{0, 100} || feature2.Input["Battery"].Command[0] != "Read" {
		t.Fatalf("unexpected battery input: %+v", feature2.Input["Battery"])
	}

	device1 := parsed.Devices[1]
	if device1.DeviceName != "Test Stroker" {
		t.Fatalf("unexpected device1: %+v", device1)
	}
	pos := device1.DeviceFeatures[0].Output["HwPositionWithDuration"]
	if pos.Value != [2]int{0, 100} || pos.Duration == nil || *pos.Duration != [2]int{0, 1000} {
		t.Fatalf("unexpected position output: %+v", pos)
	}
}

func TestDeviceListEmpty(t *testing.T) {
	data := map[string]any{"DeviceList": map[string]any{"Id": 1, "Devices": map[string]any{}}}
	msg, err := messages.ParseMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.(*messages.DeviceList).Devices) != 0 {
		t.Fatal("expected empty devices")
	}
}

func TestOutputCmdVibrateSerialize(t *testing.T) {
	msg := &messages.OutputCmd{
		BaseMessage:  messages.BaseMessage{ID: 5},
		DeviceIndex:  0,
		FeatureIndex: 0,
		Command:      map[string]any{"Vibrate": map[string]any{"Value": 10}},
	}
	result := msg.ToProtocol()
	expected := []map[string]any{{
		"OutputCmd": map[string]any{
			"Id":           5,
			"DeviceIndex":  0,
			"FeatureIndex": 0,
			"Command":      map[string]any{"Vibrate": map[string]any{"Value": 10}},
		},
	}}
	assertEqualMaps(t, result, expected)
}

func TestOutputCmdPositionWithDuration(t *testing.T) {
	msg := &messages.OutputCmd{
		BaseMessage:  messages.BaseMessage{ID: 6},
		DeviceIndex:  1,
		FeatureIndex: 0,
		Command: map[string]any{
			"HwPositionWithDuration": map[string]any{"Value": 80, "Duration": 250},
		},
	}
	result := msg.ToProtocol()
	cmd := result[0]["OutputCmd"].(map[string]any)["Command"].(map[string]any)["HwPositionWithDuration"].(map[string]any)
	if cmd["Value"] != 80 || cmd["Duration"] != 250 {
		t.Fatalf("unexpected command payload: %+v", cmd)
	}
}

func TestInputCmdSerialize(t *testing.T) {
	msg := &messages.InputCmd{
		BaseMessage:  messages.BaseMessage{ID: 7},
		DeviceIndex:  0,
		FeatureIndex: 2,
		InputType:    "Battery",
		Command:      "Read",
	}
	result := msg.ToProtocol()
	expected := []map[string]any{{
		"InputCmd": map[string]any{
			"Id":           7,
			"DeviceIndex":  0,
			"FeatureIndex": 2,
			"Type":         "Battery",
			"Command":      "Read",
		},
	}}
	assertEqualMaps(t, result, expected)
}

func TestInputReadingDeserialize(t *testing.T) {
	data := map[string]any{
		"InputReading": map[string]any{
			"Id": 5, "DeviceIndex": 0, "FeatureIndex": 2,
			"Reading": map[string]any{"Battery": map[string]any{"Value": 75}},
		},
	}
	msg, err := messages.ParseMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	parsed := msg.(*messages.InputReading)
	if parsed.Reading["Battery"].Value != 75 {
		t.Fatalf("unexpected reading: %+v", parsed.Reading)
	}
}

func TestStopCmdAllDevices(t *testing.T) {
	msg := messages.NewStopCmd(8)
	result := msg.ToProtocol()
	expected := []map[string]any{{"StopCmd": map[string]any{"Id": 8, "Inputs": true, "Outputs": true}}}
	assertEqualMaps(t, result, expected)
}

func TestStopCmdSpecificDevice(t *testing.T) {
	deviceIndex := 0
	msg := &messages.StopCmd{
		BaseMessage: messages.BaseMessage{ID: 9},
		DeviceIndex: &deviceIndex,
		Inputs:      false,
		Outputs:     true,
	}
	result := msg.ToProtocol()
	expected := []map[string]any{{"StopCmd": map[string]any{"Id": 9, "DeviceIndex": 0, "Inputs": false, "Outputs": true}}}
	assertEqualMaps(t, result, expected)
}

func TestParseMessagesArray(t *testing.T) {
	data := []map[string]any{
		{"Ok": map[string]any{"Id": 1}},
		{"ScanningFinished": map[string]any{"Id": 0}},
	}
	parsed, err := messages.ParseMessages(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(parsed))
	}
}

func TestParseMessageUnknownType(t *testing.T) {
	_, err := messages.ParseMessage(map[string]any{"UnknownType": map[string]any{"Id": 1}})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseMessageMultipleTypes(t *testing.T) {
	_, err := messages.ParseMessage(map[string]any{
		"Ok":    map[string]any{"Id": 1},
		"Error": map[string]any{"Id": 1},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestJSONRoundtrip(t *testing.T) {
	msg := messages.NewRequestServerInfo(1, "Test")
	jsonBytes, err := json.Marshal(msg.ToProtocol())
	if err != nil {
		t.Fatal(err)
	}
	var parsedData []map[string]any
	if err := json.Unmarshal(jsonBytes, &parsedData); err != nil {
		t.Fatal(err)
	}
	restored, err := messages.ParseMessage(parsedData[0])
	if err != nil {
		t.Fatal(err)
	}
	if restored.(*messages.RequestServerInfo).ClientName != "Test" {
		t.Fatalf("unexpected client name: %+v", restored)
	}
}

func assertEqualMaps(t *testing.T, got, want []map[string]any) {
	t.Helper()
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotJSON) != string(wantJSON) {
		t.Fatalf("maps differ:\ngot:  %s\nwant: %s", gotJSON, wantJSON)
	}
}
