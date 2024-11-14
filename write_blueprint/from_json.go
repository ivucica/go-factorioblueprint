package write_blueprint

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/klauspost/compress/zlib"
)

type blueprintEncoder struct {
	w io.Writer

	wroteVersion bool
}

func (b *blueprintEncoder) Write(p []byte) (n int, err error) {
	// Did we write the version byte?
	if !b.wroteVersion {
		// Write version byte.
		if x, err := b.w.Write([]byte{'0'}); x != 1 || err != nil {
			return 0, fmt.Errorf("failed to write version byte: %w", err)
		}
		b.wroteVersion = true

		// This means we did not wrap the passed writer yet.
		// Build a compressor and wrap it with encoder.
		b.w = base64.NewEncoder(base64.StdEncoding, b.w)
		b.w = zlib.NewWriter(b.w)
	}

	// Pass the rest of the data through the compression and encoding.
	return b.w.Write(p)
}

func (b *blueprintEncoder) Close() error {
	// See if our writer is a closer.
	if c, ok := b.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (b *blueprintEncoder) Flush() error {
	// See if our writer is a flusher.
	if f, ok := b.w.(Flusher); ok {
		return f.Flush()
	}
	return nil
}

type Flusher interface {
	Flush() error
}

type WriteCloseFlusher interface {
	io.WriteCloser
	Flusher
}

// AsStringWriter takes a writer expecting valid blueprint JSON and returns a
// writer that packs the JSON into a Factorio blueprint string.
//
// The passed byte data will be first compressed using zlib, then base64
// encoded. The version byte is prepended to the data.
//
// You need to finalize the data via Flush or Close if you want to read the
// zlib-encoded data correctly -- otherwise you will read only '0' (the version
// prefix rune / byte), even if you otherwise do not close or flush the writer.
//
// Close is more likely to work; please use it unless you really can't do
// without it. Encode->Decode tests did not work when Flush was used.
func AsStringWriter(w io.Writer) WriteCloseFlusher {
	// The actual wrapping is done in the blueprintEncoder's Writer.
	return &blueprintEncoder{w: w}
}
