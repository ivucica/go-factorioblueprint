package write_blueprint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"badc0de.net/pkg/factorioblueprint/read_blueprint"
	"badc0de.net/pkg/factorioblueprint/schema/blueprint_schema"
)

// Example of how to construct a valid blueprint string from a struct.
func ExampleFromStruct() {
	ptrString := func(s string) *string { return &s }
	ptrSignalIDType := func(s blueprint_schema.SignalIDType) *blueprint_schema.SignalIDType { return &s }
	ptrInt := func(i int) *int { return &i }
	var m = blueprint_schema.BlueprintSchemaJSON{
		Blueprint: &blueprint_schema.Blueprint{
			Icons: []blueprint_schema.Icon{
				{
					Signal: blueprint_schema.SignalID{
						Type: ptrSignalIDType(blueprint_schema.SignalIDTypeItem),
						Name: "transport-belt",
					},
					Index: 1,
				},
			},
			Entities: []blueprint_schema.Entity{
				{
					EntityNumber: 1,
					Name:         "transport-belt",
					Position: blueprint_schema.Position{
						X: 0,
						Y: 0,
					},
					Direction: ptrInt(6),
				},
			},
			Label:   ptrString("Blueprint"),
			Item:    "blueprint",
			Version: 281479273986304,
		},
	}

	// Create a blueprint string.
	var buf = new(bytes.Buffer)
	err := FromStruct(buf, m)
	if err != nil {
		panic(err)
	}

	// Print the blueprint.
	fmt.Print(buf.String())

	// Output:
	// 0eJx0zk1qAzEMBeB9j/HWLkx+yI+WvUYpxU5EEXhkYyslg/Hdy8Sb2XQnnqSP1xDig3MRNVADq4kJV9Bnw10K30ySgk5urJZvfcyBC2jnoH5mEKx4rTkVew8cDQ45VRlvDU/Q5LCApt6/HOSWdOCid36+mCo/6uN6/B9oS15zMZ4xmHWiTXWH6ANHED422S+X+uqxv+yO5+v+fLheTofp2PvbHwAAAP//
}

// Example of encoding an existing string into a blueprint string. It will only
// work if it's already JSON.
func ExampleAsStringWriter() {
	// Data from the blueprint string.
	const blueprint = `{
	"blueprint":
	  {
	    "entities":[{
		  "entity_number":1,
		  "name":"transport-belt",
		  "position":{"x":0,"y":0},
		  "direction":6
		  }],
		"icons":[
		  {"signal": {"type": "item", "name": "transport-belt"}, "index": 1}
		],
		"item":"blueprint",
		"label":"Blueprint",
		"version":281479273986304
	  }
}`

	// Create a blueprint string using AsStringWriter.
	//
	// (We should be able to pass a zero-value buffer to AsStringWriter, but
	// why not use a buffer with a capacity related to the blueprint string
	// length?)
	var buf = bytes.NewBuffer(make([]byte, 0, len(blueprint)*2))
	encoder := AsStringWriter(buf) // write into buf what we generate below
	if _, err := encoder.Write([]byte(blueprint)); err != nil {
		panic(err)
	}

	// Mandatory flush or close before reading the buffer.
	if err := encoder.Flush(); err != nil {
		panic(err)
	}

	// Print the blueprint.
	fmt.Print(buf.String())

	// Output:
	// 0eJxckN1q8zAMho/tqzDvcT5I2tIfHX63McZIWjEEjhNsdTQY3/twmsK2Ex08D9IrKVuDwd95jhIUZI1zuRbnwEFFhRPoLVtjXmT5CPdx4AjqmicO/cggaOxDmqeo/wb2ik3OUxKVKYAyHqC2wQJqy2ZvEvn61MeVlPdqINcp1OCVZST5DL0HuQxdZgY5iPKI5hXu/qaXxkHCjR8g1xVrzDa3dtGPi1fq+4E9CP9/4y+OaV1td+4Op8vutL+cj/v2UP9TbPkGAAD/
}

// Example of how to construct human-readable blueprint JSON string from a
// struct.
func ExampleFromStruct_humanReadable() {
	ptrString := func(s string) *string { return &s }
	ptrSignalIDType := func(s blueprint_schema.SignalIDType) *blueprint_schema.SignalIDType { return &s }
	ptrInt := func(i int) *int { return &i }
	m := blueprint_schema.BlueprintSchemaJSON{
		Blueprint: &blueprint_schema.Blueprint{
			Icons: []blueprint_schema.Icon{
				{
					Signal: blueprint_schema.SignalID{
						Type: ptrSignalIDType(blueprint_schema.SignalIDTypeItem),
						Name: "transport-belt",
					},
					Index: 1,
				},
			},
			Entities: []blueprint_schema.Entity{
				{
					EntityNumber: 1,
					Name:         "transport-belt",
					Position: blueprint_schema.Position{
						X: 0,
						Y: 0,
					},
					Direction: ptrInt(6),
				},
			},
			Label:   ptrString("Blueprint"),
			Item:    "blueprint",
			Version: 281479273986304,
		},
	}

	// Create a blueprint string.
	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		panic(err)
	}

	// Print the blueprint.
	fmt.Print(string(b))

	// Output:
	// {
	// 	"blueprint": {
	// 		"entities": [
	// 			{
	// 				"direction": 6,
	// 				"entity_number": 1,
	// 				"name": "transport-belt",
	// 				"position": {
	// 					"x": 0,
	// 					"y": 0
	// 				}
	// 			}
	// 		],
	// 		"icons": [
	// 			{
	// 				"index": 1,
	// 				"signal": {
	// 					"name": "transport-belt",
	// 					"type": "item"
	// 				}
	// 			}
	// 		],
	// 		"item": "blueprint",
	// 		"label": "Blueprint",
	// 		"version": 281479273986304
	// 	}
	// }
}

// Test that strings can be decoded losslessly.
func TestDecodeable(t *testing.T) {
	// Test a simple string.
	testInput := "some test"

	// Make output buffer.
	var buf bytes.Buffer

	// Encode.
	w := AsStringWriter(&buf)
	if _, err := w.Write([]byte(testInput)); err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	// Close.
	if err := w.Close(); err != nil {
		t.Fatalf("failed to flush: %v", err)
	}

	t.Logf("encoded: %v", buf.String())

	// Read this back.
	var r io.Reader = bytes.NewReader(buf.Bytes())

	// Decode.
	decoded, err := read_blueprint.AsJSONReader(r)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	b, err := io.ReadAll(decoded)
	if err != nil {
		t.Fatalf("failed to read all: %v", err)
	}

	// Verify b contains the same string.
	if string(b) != testInput {
		t.Fatalf("decoded string is different: %v", string(b))
	}
}
