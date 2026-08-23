//go:build tools
// +build tools

package factorioblueprint

// Declare tools we will use.
//
// Use with cat tools.go | grep -P '^\t_' | awk -F'"' '{print $2}' | xargs -tI % go install %
//
// See https://marcofranssen.nl/manage-go-tools-via-go-modules
import (
	_ "codeberg.org/emersion/go-jsonschema" // @v0.0.0-20251116133759-a828df140a57"
	_ "github.com/atombender/go-jsonschema" // @v0.15
)
