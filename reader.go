package b64reader

import (
	"bytes"
	"encoding/base64"
	"io"
)

type At struct {
	io.ReaderAt
}

func (at At) ReadAt(p []byte, off int64) (int, error) {
	enc := base64.StdEncoding

	subStart := off / 3 * 4
	subSkip := off * 4 % 3
	subLen := int(subSkip) + len(p)

	buf := make([]byte, enc.EncodedLen(subLen))

	if n, err := at.ReaderAt.ReadAt(buf, subStart); err != nil && err != io.EOF {
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
