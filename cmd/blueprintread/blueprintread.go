// blueprintread reads a b64-encoded zlib-compressed blueprint string from a
// file, which is in JSON format at that point, then tries to read it into a
// schema, and print it out in some form.
package main // badc0de.net/pkg/factorioblueprint/cmd/blueprintread

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"badc0de.net/pkg/factorioblueprint/asciiart_blueprint"
	"badc0de.net/pkg/factorioblueprint/read_blueprint"

	"gopkg.in/yaml.v3"
)

var (
	file   = flag.String("file", "", "The file to read the blueprint from. If empty, uses stdin.")
	format = flag.String("fmt", "json", "Format. raw_json (no processing after decompression), json (default, pretty print JSON), yaml, asciiart (experimental and halfbroken).")
)

func init() {
	flag.Parse()
}

func main() {
	r := os.Stdin
	if *file != "" {
		var err error
		r, err = os.Open(*file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
			os.Exit(1)
		}
	}
	defer r.Close()

	decompressed, err := read_blueprint.AsJSONReader(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decompress JSON: %v\n", err)
		os.Exit(1)
	}

	// Print out raw JSON for debugging.
	if *format == "raw_json" {
		rawJSON, err := ioutil.ReadAll(decompressed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read raw JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s", rawJSON)
		return
	}

	// Decode JSON.
	//var m = make(map[string]interface{})
	m, err := read_blueprint.AsStruct(decompressed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode JSON: %v\n", err)
		os.Exit(1)
	}

	switch *format {
	case "raw_json":
		panic("unreachable")
	case "json":
		// Print out marshalled prettified JSON.
		if b, err := json.MarshalIndent(m, "", "  "); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("%s\n", b)
		}
	case "yaml":
		// Print out YAML.
		if b, err := yaml.Marshal(m); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal YAML: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("%s\n", b)
		}
	case "asciiart":
		// Print out ASCII art of the tilemap. Just use 1x1 for now.
		r := asciiart_blueprint.NewReader(m.Blueprint, 1, 1)
		// Copy to stdout.
		if _, err := io.Copy(os.Stdout, r); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate ASCII art: %v\n", err)
			os.Exit(1)
		}
		fmt.Println() // no newline from generated asciiart, so add one
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %v\n", *format)
		os.Exit(1)
	}

}
