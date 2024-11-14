#!/bin/bash

if [[ -z "$1" ]] ; then
  echo "Usage: $0 <blueprint.txt>"
  exit 1
fi

# Skip the first character, which is a version number.
# Decode the base64-encoded string.
# Decompress the gzip-compressed string.
# Pretty-print the JSON.
cut -c2- "$1" | base64 -d | pigz -d | jq .
