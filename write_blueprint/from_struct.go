package write_blueprint

import (
	"encoding/json"
	"fmt"
	"io"

	"badc0de.net/pkg/factorioblueprint/schema/blueprint_schema"
)

// FromStruct takes a blueprint schema and an io.Writer and writes the blueprint
// to the writer.
func FromStruct(w io.Writer, m blueprint_schema.BlueprintSchemaJSON) error {
	// Create a blueprint encoder.
	encoder := AsStringWriter(w)

	// Write the JSON.
	if err := json.NewEncoder(encoder).Encode(m); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Flush or close in order for zlib to encode the block. Flush should be
	// enough.
	if err := encoder.Flush(); err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	return nil
}
