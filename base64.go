package tio

import (
	"bytes"
	"cmp"
	"encoding/base64"
	"io"
)

// Base64ReaderAt is a wrapper around an existing [io.ReaderAt] that assumes the underlying data is base64-encoded and will transparently decode it during reads.  See also https://github.com/jonjohnsonjr/targz
type Base64ReaderAt struct {
	// R is the [io.ReaderAt] to ReadAt base64 data from
	R   io.ReaderAt

	// Enc is the [base64.Encoding] to use for decoding -- if unset, will default to [base64.StdEncoding] automatically
	Enc *base64.Encoding
}

func (at Base64ReaderAt) ReadAt(p []byte, off int64) (int, error) {
	enc := cmp.Or(at.Enc, base64.StdEncoding)

	subStart := off / 3 * 4
	subSkip := off * 4 % 3
	subLen := int(subSkip) + len(p)

	buf := make([]byte, enc.EncodedLen(subLen))

	if n, err := at.R.ReadAt(buf, subStart); err != nil && err != io.EOF {
		return 0, err
	} else if n < len(buf) {
		buf = buf[:n]
	}

	r := base64.NewDecoder(enc, bytes.NewReader(buf))

	if _, err := io.CopyN(io.Discard, r, subSkip); err != nil {
		return 0, err
	}

	// I really wish the stdlib had something like "io.ReadAtMost" (some combination of io.LimitedReader and io.ReadAll but into a provided buffer)
	// I can emulate it with io.ReadAtLeast and converting io.ErrUnexpectedEOF into io.EOF, but it's hanky
	n, err := io.ReadAtLeast(r, p, len(p))
	if err == io.ErrUnexpectedEOF || n < len(p) {
		err = io.EOF
	}
	return n, err
}
