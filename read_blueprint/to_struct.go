package read_blueprint

import (
	"encoding/json"
	"io"

	"badc0de.net/pkg/factorioblueprint/schema/blueprint_schema"
)

func AsStruct(decompressed io.Reader) (m blueprint_schema.BlueprintSchemaJSON, err error) {
	d := json.NewDecoder(decompressed)
	if err := d.Decode(&m); err != nil {
		return m, err
	}
	return m, nil
}
