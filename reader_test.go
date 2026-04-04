package b64reader

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestReaderAt(t *testing.T) {
	at := At{bytes.NewReader([]byte("AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8="))} // 0x00, 0x01, 0x02, ..., 0x10, ..., 0x1F (32 bytes)
	for _, td := range []struct {
		p   []byte
		off int64
	}{
		{[]byte{0}, 0},
		{[]byte{1}, 1},
		{[]byte{2}, 2},
		{[]byte{3}, 3},
		{[]byte{4}, 4},
		{[]byte{8, 9}, 8},
		{[]byte{30, 31}, 30},
		{nil, 32},
	} {
		t.Run(fmt.Sprintf("%d-%d", td.off, len(td.p)), func(t *testing.T) {
			buf := make([]byte, len(td.p))
			if td.p == nil {
				buf = make([]byte, 30)
			}
			n, err := at.ReadAt(buf, td.off)
			if err != nil && err != io.EOF {
				t.Fatal(err)
			}
			if n != len(td.p) {
				t.Fatalf("%d != %d", n, len(td.p))
			}
			if !bytes.Equal(buf[:n], td.p) {
				t.Fatalf("%v != %v", buf, td.p)
			}
		})
	}
}
