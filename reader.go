package b64reader

import (
	"encoding/base64"
	"io"
)

type At struct {
	RA io.ReaderAt
}

func (at At) ReadAt(p []byte, off int64) (n int, err error) {
	enc := base64.StdEncoding
	subStart := off / 3 * 4
	subSkip := int(off * 4 % 3)
	subLen := enc.EncodedLen(subSkip + len(p))
	buf := make([]byte, subLen)
	if n, err := at.RA.ReadAt(buf, int64(subStart)); err != nil {
		return 0, err
	} else if n < len(buf) {
		buf = buf[:n]
	}
	// TODO bit math instead?
	p2 := make([]byte, enc.DecodedLen(len(buf)))
	if n, err := enc.Decode(p2, buf); err != nil {
		if _, ok := err.(base64.CorruptInputError); ok && n == len(p) {
			err = nil
		}
	}
	p2 = p2[subSkip : subSkip+len(p)]
	copy(p, p2) // TODO something clever to avoid this copy?? (bit trimming above?)
	return len(p2), err
}
