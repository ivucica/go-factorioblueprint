package asciiart_blueprint

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"badc0de.net/pkg/factorioblueprint/schema/blueprint_schema"
)

// setupSimpleBlueprint creates a simple blueprint struct for testing.
//
// A blueprint with a transport-belt, facing east, and an inserter loading from
// the belt into a chest, which are both above the belt.
func setupSimpleBlueprint() *blueprint_schema.BlueprintSchemaJSON {
	ptrString := func(s string) *string { return &s }
	ptrSignalIDType := func(s blueprint_schema.SignalIDType) *blueprint_schema.SignalIDType { return &s }
	ptrInt := func(i int) *int { return &i }
	m := &blueprint_schema.BlueprintSchemaJSON{
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
						X: -0.5,
						Y: 1.5,
					},
					Direction: ptrInt(6),
				},
				{
					EntityNumber: 2,
					Name:         "inserter",
					Position: blueprint_schema.Position{
						X: -0.5,
						Y: 0.5,
					},
					Direction: ptrInt(4), // taking from down, inserting to up
				},
				{
					EntityNumber: 3,
					Name:         "stone-furnace",
					Position: blueprint_schema.Position{
						X: 0,
						Y: -1,
					},
				},
				{
					EntityNumber: 4,
					Name:         "transport-belt",
					Position: blueprint_schema.Position{
						X: -1.5,
						Y: 1.5,
					},
					Direction: ptrInt(6),
				},
			},
			Label:   ptrString("Blueprint"),
			Item:    "blueprint",
			Version: 281479273986304,
		},
	}
	return m
}

// ExampleReader demonstrates how to use the asciiart blueprint reader with the
// blueprint already read into the struct.
//
// Because this example requests minimal size (1x1) for each tile, there is
// enough space only to represent the entities ("tiles" or "entities").
func ExampleReader_has1x1() {
	// Create a simple blueprint.
	m := setupSimpleBlueprint()

	// Create a reader for the blueprint.
	//
	// Note how the JSON blueprint, and not the entire file, is passed in. This
	// is because we do not support printing out the entire book.
	r := NewReader(m.Blueprint, 1, 1)

	// Print the blueprint to stdout.
	// io.Copy(os.Stdout, r)
	// However, for purposes of the testing framework in Go, we can't depend on
	// whitespace. Place it into a temp buffer and replace with dots.
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		panic(err) // do something nicer instead
	}
	out := strings.Replace(buf.String(), " ", ".", -1)
	fmt.Print(out)

	// BUG: there is extra spacing in the sample output, due to misunderstanding
	// about the appearance of 0.5 indicating that tiles are 0.5x0.5. Uses of
	// 0.5 should be removed.

	// TODO: the 2x2 furnace is actually in top right, not in top left, of
	// its expected position. Does Factorio actually store center of the object
	// as its position? Maybe instead of *2 we need to +=0.5 to everything?

	// See also comments in ExampleReader_2x2.

	// Output:
	// ...k
	// ....
	// ....
	// ..I.
	// ....
	// v.v.
	// ---
	// [I]:.inserter
	// [k]:.stone-furnace
	// [v]:.transport-belt
}

// ExampleReader_2x2 demonstrates how to use the asciiart blueprint reader with
// a 2x2 size.
func ExampleReader_has2x2() {
	// Create a simple blueprint.
	m := setupSimpleBlueprint()

	// Create a reader for the blueprint.
	//
	// Note how the JSON blueprint, and not the entire file, is passed in. This
	// is because we do not support printing out the entire book.
	r := NewReader(m.Blueprint, 2, 2)

	// Print the blueprint to stdout.
	// io.Copy(os.Stdout, r)
	// However, for purposes of the testing framework in Go, we can't depend on
	// whitespace. Place it into a temp buffer and replace with dots.
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		panic(err) // do something nicer instead
	}
	out := strings.Replace(buf.String(), " ", ".", -1)
	fmt.Print(out)

	// BUG: same bugs and notes apply as with ExampleReader (extra whitespace).
	//
	// BUG: the furnace is drawn twice, once clipping outside the tilemap.
	//
	// BUG: the directions are centered differently than the entities themselves
	// which is unintended.
	//
	// Note that the blueprint is actually something like this:
	// .kk
	// .kk
	// .I.
	// vv.
	//
	// Or on this 2x2:
	// ..kkkk
	// ..kkkk
	// ..kkkk
	// ..kkkk
	// ..II..
	// ..II..
	// vvvv..
	// ......
	//
	// then given how we actually do 2x2 to draw directions, we should expect
	// something resembling this (especially since we don't have knowledge that
	// the furnace is 2x2):
	// .....k
	// ...... (no direction)
	// ......
	// ......
	// ...I..
	// ...V..
	// .v.v..
	// .<.<..

	// Output:
	// ......k.k
	// ........
	// ........
	// ........
	// ........
	// ........
	// ....I...
	// ...V....
	// ........
	// ........
	// .v..v...
	// <..<....
	// ---
	// [I]:.inserter
	// [k]:.stone-furnace
	// [v]:.transport-belt
}

// disabled_ExampleReader_has3x3 shows how the reader prints a blueprint using
// 3x3 tile size.
//
// BUG: The test is broken and outputs nothing.
func disabled_ExampleReader_has3x3() {
	// Create a simple blueprint.
	m := setupSimpleBlueprint()

	// Create a reader for the blueprint.
	//
	// Note how the JSON blueprint, and not the entire file, is passed in. This
	// is because we do not support printing out the entire book.
	r := NewReader(m.Blueprint, 3, 3)

	// Print the blueprint to stdout.
	// io.Copy(os.Stdout, r)
	// However, for purposes of the testing framework in Go, we can't depend on
	// whitespace. Place it into a temp buffer and replace with dots.
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		panic(err) // do something nicer instead
	}

	out := strings.Replace(buf.String(), " ", ".", -1)
	fmt.Print(out)

	// BUG: same bugs and notes apply as with ExampleReader (extra whitespace).
	//
	// BUG: the furnace is drawn incorrectly
	//
	// BUG: the directions are centered differently than the entities themselves
	// which is unintended.

	// Output:
	// nothing here.
}

// TestReader_Size tests that the reader correctly computes the size of the
// blueprint.
func TestReader_Size(t *testing.T) {
	tcs := []struct {
		name         string
		blueprint    *blueprint_schema.BlueprintSchemaJSON
		wantW, wantH float64 // in factorio coordinates (so precision of 0.5)
	}{
		{
			name:      "SimpleBlueprint",
			blueprint: setupSimpleBlueprint(),
			wantW:     2.5, // we don't know that the furnace is 2x2 otherwise we'd expect 3.5
			wantH:     3.5,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader(tc.blueprint.Blueprint, 1, 1)
			gotW, gotH := r.Size()
			if gotW != tc.wantW || gotH != tc.wantH {
				t.Errorf("expected size %gx%g, got %gx%g", tc.wantW, tc.wantH, gotW, gotH)
			}
		})
	}
}

// TestReader_box tests that the reader correctly computes the box of the
// blueprint (min xy, max xy).
func TestReader_box(t *testing.T) {
	tcs := []struct {
		name      string
		blueprint *blueprint_schema.BlueprintSchemaJSON

		wantMinX, wantMinY, wantMaxX, wantMaxY float64 // in factorio coordinates (so precision of 0.5)
	}{
		{
			name:      "SimpleBlueprint",
			blueprint: setupSimpleBlueprint(),
			wantMinX:  -1.5,
			wantMinY:  -1,
			wantMaxX:  1,   // max is 0, we add 1; furnace is 2x2 so if we supported other data, we'd have 2
			wantMaxY:  2.5, // max is 1.5, we add 1
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader(tc.blueprint.Blueprint, 1, 1)
			gotMinX, gotMinY, gotMaxX, gotMaxY := r.box()
			if gotMinX != tc.wantMinX || gotMinY != tc.wantMinY || gotMaxX != tc.wantMaxX || gotMaxY != tc.wantMaxY {
				t.Errorf("expected box %g,%g - %g,%g, got %g,%g - %g,%g", tc.wantMinX, tc.wantMinY, tc.wantMaxX, tc.wantMaxY, gotMinX, gotMinY, gotMaxX, gotMaxY)
			}
		})
	}
}
