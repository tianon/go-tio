package b64reader

import (
	"encoding/base64"
	"io"
)

type At struct {
	RA io.ReaderAt
}

func (at At) ReadAt(p []byte, off int64) (int, error) {
	enc := base64.StdEncoding

	subStart := off / 3 * 4
	subSkip := int(off * 4 % 3)
	subLen := subSkip + len(p)

	buf := make([]byte, enc.EncodedLen(subLen))

	if n, err := at.RA.ReadAt(buf, int64(subStart)); err != nil && err != io.EOF {
		return 0, err
	} else if n < len(buf) {
		buf = buf[:n]
	}
	// TODO clever bit math / shifting / handling?
	p2 := make([]byte, enc.DecodedLen(len(buf)))
	if n, err := enc.Decode(p2, buf); err != nil {
		if _, isCorruptInputError := err.(base64.CorruptInputError); isCorruptInputError && n == subLen {
			p2 = p2[:n]
		} else {
			return 0, err
		}
	} else {
		p2 = p2[:n]
	}
	p2 = p2[subSkip:min(subLen, len(p2))]
	copy(p, p2) // TODO something clever to avoid this copy?? (bit trimming above?)
	if len(p2) < len(p) {
		return len(p2), io.EOF
	}
	return len(p2), nil
}
