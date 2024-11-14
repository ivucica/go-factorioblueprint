//go:build tools
// +build tools

package factorioblueprint

// Declare tools we will use.
//
// Use with cat tools.go | grep _ | awk -F'"' '{print $2}' | xargs -tI % go install %
//
// See https://marcofranssen.nl/manage-go-tools-via-go-modules
import (
	_ "github.com/atombender/go-jsonschema" // @v0.14
)
