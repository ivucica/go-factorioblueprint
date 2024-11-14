package read_blueprint

import (
	_ "embed"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"badc0de.net/pkg/factorioblueprint/schema/blueprint_schema"
)

// SimpleTxt is a simple blueprint with some belts and inserters.
//
//go:embed simple.txt
var SimpleTxt string

// SimpleJSON is the JSON representation of SimpleTxt generated using tojson.sh.
//
//go:embed simple.json
var SimpleJSON string

// All the test data.
var tests = []struct {
	filename string
	content  string
}{
	{"simple.txt", SimpleTxt},
	{"simple.json", SimpleJSON},
}

func TestAsJSONReader(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			r := strings.NewReader(tc.content)
			decompressed, err := AsJSONReader(r)
			if err != nil {
				t.Fatalf("Failed to decompress JSON: %v", err)
			}

			// Read all.
			data, err := ioutil.ReadAll(decompressed)
			if err != nil {
				t.Fatalf("Failed to read decompressed JSON: %v", err)
			}

			// No error.
			// Read JSON into a map.
			var m = make(map[string]interface{})
			if err := json.Unmarshal(data, &m); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			// Check some values.
			// Root should be 'blueprint' containing 'entities', 'icons', 'item' and
			// 'label'.
			// item should be "blueprint".
			if bp := m["blueprint"].(map[string]interface{})["item"]; bp != "blueprint" {
				t.Fatalf("Bad data: want = 'blueprint' got '%v'", bp)
			}
			// label should be "Blueprint".
			if bp := m["blueprint"].(map[string]interface{})["label"]; bp != "Blueprint" {
				t.Fatalf("Bad data: want = 'Blueprint' got '%v'", bp)
			}
			// Icon at index 0 should be "transport-belt" under the value "signal".
			bp := m["blueprint"].(map[string]interface{})
			icons := bp["icons"].([]interface{})
			signal := icons[0].(map[string]interface{})["signal"]
			if signalName := signal.(map[string]interface{})["name"]; signalName != "transport-belt" {
				t.Fatalf("Bad data: want = 'transport-belt' got '%v'", signalName)
			}
			// signalType should be "item".
			if signalType := signal.(map[string]interface{})["type"]; signalType != "item" {
				t.Fatalf("Bad data: want = 'item' got '%v'", signalType)
			}
			// That same icon should have index value 1, a number.
			index := icons[0].(map[string]interface{})["index"]
			if index != float64(1) {
				t.Fatalf("Bad data: want = 1 got '%v'", index)
			}
		})
	}
}

func TestAsStruct(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			r := strings.NewReader(tc.content)
			decompressed, err := AsJSONReader(r)
			if err != nil {
				t.Fatalf("Failed to decompress JSON: %v", err)
			}

			// Read all.
			data, err := ioutil.ReadAll(decompressed)
			if err != nil {
				t.Fatalf("Failed to read decompressed JSON: %v", err)
			}

			// No error.
			// Read JSON into a struct.
			var m blueprint_schema.BlueprintSchemaJSON
			if err := json.Unmarshal(data, &m); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			// Check some values.
			// Root should be 'blueprint' containing 'entities', 'icons', 'item' and
			// 'label'.
			// item should be "blueprint".
			if m.Blueprint.Item != "blueprint" {
				t.Fatalf("Bad data: want = 'blueprint' got '%v'", m.Blueprint.Item)
			}
			// label should be "Blueprint".
			if m.Blueprint.Label != nil && *m.Blueprint.Label != "Blueprint" {
				t.Fatalf("Bad data: want = 'Blueprint' got '%v'", m.Blueprint.Label)
			}
			// Icon at index 0 should be "transport-belt" under the value "signal".
			signalName := m.Blueprint.Icons[0].Signal.Name
			if signalName != "transport-belt" {
				t.Fatalf("Bad data: want = 'transport-belt' got '%v'", signalName)
			}
			// signalType should be "item".
			signalType := m.Blueprint.Icons[0].Signal.Type
			if signalType != nil && *signalType != "item" {
				t.Fatalf("Bad data: want = 'item' got '%v'", signalType)
			}
			// That same icon should have index value 1, a number.
			index := m.Blueprint.Icons[0].Index
			if index != 1 {
				t.Fatalf("Bad data: want = 1 got '%v'", index)
			}
		})
	}
}
