package eport

import (
	"bytes"
	"io"
	"testing"
)

type bs []byte

type rfun func() ([]byte, error)

func newRW(from []byte) (*bytes.Reader, *bytes.Buffer) {
	return bytes.NewReader(from), new(bytes.Buffer)
}

func testRead(
	t *testing.T,
	f rfun,
	sizes []int,
	datas []bs,
	errs []error) {
	var data []byte
	var err error

	if len(sizes) != len(datas) || len(sizes) != len(errs) {
		t.Fatalf("wrong lengths: sizes=%d, datas=%d, errs=%d",
			len(sizes), len(datas), len(errs))
	}

	for i, size := range sizes {
		data, err = f()
		if err != errs[i] {
			t.Errorf("failed on %d: got error %v, should be %v", i, err, errs[i])
		}

		if size != len(data) {
			t.Errorf("failed on %d: got data (%s) size %d, should be %d",
				i, data, len(data), size)
		}

		if bytes.Compare(datas[i], data) != 0 {
			t.Errorf("failed on %d: got data %v, should be %v", i, data, datas[i])
		}
	}

	data, err = f()

	if err != io.EOF {
		t.Errorf("last error value must be io.EOF, got %v instead", err)
	} else if len(data) != 0 {
		t.Errorf("non-empty data (%v) after EOF", data)
	}
}

func testWrite(t *testing.T, p Port, w *bytes.Buffer, sizes []int, outs []bs, final bs) {
	for i, out := range outs {
		n, err := p.Write(out)
		if n != sizes[i] {
			t.Errorf("failed on %d: expected size %d, got %d", i, sizes[i], n)
		}

		if err != nil {
			t.Errorf("failed on %d: %#v", err)
		}
	}

	if bytes.Compare(w.Bytes(), final) != 0 {
		t.Errorf("expected %#v, got %#v", final, w.Bytes())
	}
}
