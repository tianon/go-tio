package tio_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"testing"

	"go.tianon.xyz/tio"
)

func makeB64(size int) (b []byte, b64 []byte) {
	b = make([]byte, size)
	for i := range len(b) {
		b[i] = byte(i)
	}
	b64 = []byte(base64.StdEncoding.EncodeToString(b))
	return b, b64
}

func BenchmarkBase64ReaderAt(b *testing.B) {
	for ex := range 13 {
		kib := 1 << ex
		b.Run(fmt.Sprintf("%dKiB", kib), func(b *testing.B) {
			ogb, b64 := makeB64(kib * 1024)
			at := tio.Base64ReaderAt{R: bytes.NewReader(b64)}
			buf := make([]byte, len(ogb)/4)
			for b.Loop() {
				at.ReadAt(buf, int64(len(ogb)/2))
			}
		})
	}
}

func TestBase64ReaderAt(t *testing.T) {
	b, b64 := makeB64(64)
	buf := make([]byte, len(b)+5)
	for extraLayers := range 10 {
		var at io.ReaderAt
		b64d := b64
		for range extraLayers {
			b64d = []byte(base64.StdEncoding.EncodeToString(b64d))
		}
		at = tio.Base64ReaderAt{R: bytes.NewReader(b64d)}
		for range extraLayers {
			at = tio.Base64ReaderAt{R: at}
		}
		t.Run(fmt.Sprintf("%d layers", 1+extraLayers), func(t *testing.T) {
			for off := range len(b) + 1 { // +1 to make sure we test reading off the end
				for size := range len(b) - off + 5 { // +5 to make sure we test reading way off the end at least one full base64 "chunk" or whatever it's called
					exp := b[off:min(off+size, len(b))]
					t.Run(fmt.Sprintf("%d+%d", off, size), func(t *testing.T) {
						n, err := at.ReadAt(buf[:size], int64(off))
						if err != nil && err != io.EOF {
							t.Fatal(err)
						} else if n != size && err != io.EOF {
							t.Error("error should be io.EOF")
						} else if n == size && off+size == len(b) && err != io.EOF {
							t.Log("error could be io.EOF but isn't")
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
		})
	}
}

func TestBase64ReaderAtInvalidBase64(t *testing.T) {
	at := tio.Base64ReaderAt{R: bytes.NewReader([]byte("!!!!!!!!!! this is not valid base64! lmao"))}
	buf := make([]byte, 5)
	n, err := at.ReadAt(buf, 5)
	if err == nil {
		t.Errorf("expected error, read %d bytes instead: %+v", n, buf[:n])
	} else if err == io.EOF {
		t.Error("expected interesting error, got EOF instead")
	}
}

type readerAtError string

func (str readerAtError) ReadAt(_ []byte, _ int64) (int, error) {
	return 0, errors.New(string(str))
}

func TestBase64ReaderAtErroringReaderAt(t *testing.T) {
	at := tio.Base64ReaderAt{R: readerAtError("always an error")}
	buf := make([]byte, 1)
	n, err := at.ReadAt(buf, 1)
	if err == nil {
		t.Errorf("expected error, read %d bytes instead: %+v", n, buf[:n])
	} else if err == io.EOF {
		t.Error("expected interesting error, got EOF instead")
	}
}
