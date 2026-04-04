package b64reader

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"testing"
)

func TestReaderAt(t *testing.T) {
	b := make([]byte, 64)
	for i := range len(b) {
		b[i] = byte(i)
	}
	at := At{bytes.NewReader([]byte(base64.StdEncoding.EncodeToString(b)))}
	buf := make([]byte, len(b)+5)
	for off := range len(b) + 1 { // +1 to make sure we test reading off the end
		for size := range len(b) - off + 5 { // +5 to make sure we test reading way off the end at least one full base64 "chunk" or whatever it's called
			exp := b[off:min(off+size, len(b))]
			t.Run(fmt.Sprintf("%d+%d", off, size), func(t *testing.T) {
				n, err := at.ReadAt(buf[:size], int64(off))
				if err != nil && err != io.EOF {
					t.Fatal(err)
				} else if n != size && err != io.EOF {
					t.Error("error should be io.EOF")
				}
				if n != len(exp) {
					t.Errorf("%d != %d", n, len(exp))
				}
				if !bytes.Equal(buf[:n], exp) {
					t.Errorf("%v != %v", buf[:n], exp)
				}
			})
		}
	}
}
