#!/bin/bash
if [[ -z "${GOPATH}" ]] ; then
  echo 'must be in GOPATH'
  exit 1
fi

if [[ ! -e "${GOPATH}"/bin/go-jsonschema ]] ; then
  # n.b. Prefer using tools.go.
  # go get -v github.com/atombender/go-jsonschema/...
  go install -v github.com/atombender/go-jsonschema@v0.15
fi

if realpath blueprint.schema.json 2>&1 > /dev/null ; then
  echo "Using default blueprint schema" > /dev/stderr
  SCHEMA="$(realpath blueprint.schema.json)"
  # Temporary: Avoid use of different name, we want BlueprintSchemaJSON, which is derived from the file name.
  # Note that --resolve-extension json would remove JSON from the type name.
  SCHEMA="$(dirname "${SCHEMA}")"/blueprint.schema.json
elif realpath blueprint.ext-2.0.0.schema.json 2>&1 > /dev/null ; then
  echo "Using factorio-blueprint-schema submodule's schema"
  echo "(This requires more work: the 'blueprint' schema is not usable as-is (use of object where converter does not support it), and the root one 'blueprintable' uses refs in a way that results in empty file -- and resolve-extension does not seem to work)"
  SCHEMA="$(realpath blueprint.ext-2.0.0.schema.json)"
else
  echo "Using internal schema; check if you checked out submodules (e.g. factorio-blueprint-schemas/schemas/2.0.0/blueprintable.json should exist) > /dev/stderr"
  SCHEMA="$(realpath blueprint.internal.schema.json)"
fi

(
  # Entering directory so that relative paths in $ref work.
  cd "$(dirname "${SCHEMA}")"
  pwd
  go-jsonschema \
    --capitalization ID,JSON \
    -e \
    -p badc0de.net/pkg/factorioblueprint/schema/blueprint_schema \
    "${SCHEMA}" \
    -o ${GOPATH}/src/badc0de.net/pkg/factorioblueprint/schema/blueprint_schema/blueprint.schema.json.go
)
