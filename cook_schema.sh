#!/bin/bash
if [[ -z "${GOPATH}" ]] ; then
  echo 'must be in GOPATH'
  exit 1
fi

if [[ ! -e "${GOPATH}"/bin/go-jsonschema ]] ; then
  # n.b. Prefer using tools.go.
  # go get -v github.com/atombender/go-jsonschema/...
  go install -v github.com/atombender/go-jsonschema@v0.14
fi

go-jsonschema --capitalization ID,JSON -e -p badc0de.net/pkg/factorioblueprint/schema/blueprint_schema blueprint.schema.json -o ${GOPATH}/src/badc0de.net/pkg/factorioblueprint/schema/blueprint_schema/blueprint.schema.json.go
