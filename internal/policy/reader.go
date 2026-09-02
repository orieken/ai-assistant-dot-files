package policy

import (
	"bytes"
	"io"
)

// newReader wraps raw bytes for the streaming decoder, which is what
// enables KnownFields strictness.
func newReader(raw []byte) io.Reader { return bytes.NewReader(raw) }
