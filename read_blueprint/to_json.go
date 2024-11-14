package read_blueprint

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/klauspost/compress/zlib" // we could use compress/zlib from stdlib
)

func AsJSONReader(r io.Reader) (io.Reader, error) {
	// Read first byte.
	var version = make([]byte, 1)
	if n, err := r.Read(version); err != nil {
		return nil, fmt.Errorf("failed to read version byte: %w", err)
	} else if n != 1 {
		return nil, fmt.Errorf("failed to read version byte: read %v bytes", n)
	}

	// Assert version byte is rune 0 or rune {. If it's {, assume it's JSON.
	switch v := version[0]; v {
	case '{':
		// Return a new reader with the version byte prepended.
		r2 := io.MultiReader(bytes.NewReader(version), r)
		return r2, nil
	case '0':
		// Continue.
	default:
		return nil, fmt.Errorf("version byte: got = '%c', want = '{' or '0'", v)
	}

	// Decode b64 using standard encoding and io.Reader.
	decoder := base64.NewDecoder(base64.StdEncoding, r)

	// Decompress zlib.
	decompressed, err := zlib.NewReader(decoder)
	if err != nil {
		return nil, fmt.Errorf("failed to start decompressing with zlib: %w", err)
	}

	return decompressed, nil
}
