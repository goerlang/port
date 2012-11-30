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

func TestLineRead(t *testing.T) {
	in := bs("1234567890123\n\nфыва\r\n\tолдж\n1")
	sizes := []int{0, 1, 10, 10, 1, 0}
	datas := []bs{
		bs{}, bs("\n"), bs("фыва\r\n"), bs("\tолдж\n"), bs("1"), bs{},
	}
	errors := []error{ErrTooBig, nil, nil, nil, io.EOF, io.EOF}

	r, w := newRW(in)
	p, _ := Line(r, w)
	out := make([]byte, 10)
	f := func() ([]byte, error) {
		n, err := p.Read(out)
		return out[:n], err
	}
	testRead(t, f, sizes, datas, errors)
}

func TestLineReadOne(t *testing.T) {
	in := bs("1234567890\n\nфыва\r\n\tолдж\n1")
	sizes := []int{11, 1, 10, 10, 1, 0}
	datas := []bs{
		bs("1234567890\n"), bs("\n"), bs("фыва\r\n"), bs("\tолдж\n"), bs("1"), bs{},
	}
	errors := []error{nil, nil, nil, nil, io.EOF, io.EOF}

	r, w := newRW(in)
	p, _ := Line(r, w)
	f := func() ([]byte, error) { return p.ReadOne() }
	testRead(t, f, sizes, datas, errors)
}

func TestPacket1Read(t *testing.T) {
	in := bs{
		1, 'a',
		2, 'a', 'b',
		0,
		0,
		5, '1', '2', '3', '4', '5',
	}
	sizes := []int{1, 2, 0, 0, 0, 0}
	datas := []bs{bs("a"), bs("ab"), bs{}, bs{}, bs{}, bs{}}
	errors := []error{nil, nil, nil, nil, ErrTooBig, io.EOF}

	r, w := newRW(in)
	p, _ := Packet(r, w, 1)
	out := make([]byte, 3)
	f := func() ([]byte, error) {
		n, err := p.Read(out)
		return out[:n], err
	}
	testRead(t, f, sizes, datas, errors)
}

func TestPacket1ReadOne(t *testing.T) {
	in := bs{
		1, 'a',
		2, 'a', 'b',
		0,
		0,
		5, '1', '2', '3', '4', '5',
	}
	sizes := []int{1, 2, 0, 0, 5, 0}
	datas := []bs{bs("a"), bs("ab"), bs{}, bs{}, bs("12345"), bs{}}
	errors := []error{nil, nil, nil, nil, nil, io.EOF}

	r, w := newRW(in)
	p, _ := Packet(r, w, 1)
	f := func() ([]byte, error) { return p.ReadOne() }
	testRead(t, f, sizes, datas, errors)
}

func TestPacket2Read(t *testing.T) {
	in := bs{
		0, 1, 'a',
		0, 2, 'a', 'b',
		0, 0,
		0, 5, '1', '2', '3', '4', '5',
	}
	sizes := []int{1, 2, 0, 0, 0}
	datas := []bs{bs("a"), bs("ab"), bs{}, bs{}, bs{}}
	errors := []error{nil, nil, nil, ErrTooBig, io.EOF}

	r, w := newRW(in)
	p, _ := Packet(r, w, 2)
	out := make([]byte, 3)
	f := func() ([]byte, error) {
		n, err := p.Read(out)
		return out[:n], err
	}
	testRead(t, f, sizes, datas, errors)
}

func TestPacket2ReadOne(t *testing.T) {
	in := bs{
		0, 1, 'a',
		0, 2, 'a', 'b',
		0, 0,
		0, 5, '1', '2', '3', '4', '5',
	}
	sizes := []int{1, 2, 0, 5, 0}
	datas := []bs{bs("a"), bs("ab"), bs{}, bs("12345"), bs{}}
	errors := []error{nil, nil, nil, nil, io.EOF}

	r, w := newRW(in)
	p, _ := Packet(r, w, 2)
	f := func() ([]byte, error) { return p.ReadOne() }
	testRead(t, f, sizes, datas, errors)
}

func TestPacket4Read(t *testing.T) {
	in := bs{
		0, 0, 0, 1, 'a',
		0, 0, 0, 2, 'a', 'b',
		0, 0, 0, 5, '1', '2', '3', '4', '5',
	}
	sizes := []int{1, 2, 0, 0}
	datas := []bs{bs("a"), bs("ab"), bs{}, bs{}}
	errors := []error{nil, nil, ErrTooBig, io.EOF}

	r, w := newRW(in)
	p, _ := Packet(r, w, 4)
	out := make([]byte, 3)
	f := func() ([]byte, error) {
		n, err := p.Read(out)
		return out[:n], err
	}
	testRead(t, f, sizes, datas, errors)
}

func TestPacket4ReadOne(t *testing.T) {
	in := bs{
		0, 0, 0, 1, 'a',
		0, 0, 0, 2, 'a', 'b',
		0, 0, 0, 5, '1', '2', '3', '4', '5',
	}
	sizes := []int{1, 2, 5, 0}
	datas := []bs{bs("a"), bs("ab"), bs("12345"), bs{}}
	errors := []error{nil, nil, nil, io.EOF}

	r, w := newRW(in)
	p, _ := Packet(r, w, 4)
	f := func() ([]byte, error) { return p.ReadOne() }
	testRead(t, f, sizes, datas, errors)
}

func TestStreamRead(t *testing.T) {
	in := bs("12345\n\n\x00\r")
	sizes := []int{6, 3, 0}
	datas := []bs{
		bs("12345\n"), bs("\n\x00\r"), bs{},
	}
	errors := []error{nil, nil, io.EOF}

	r, w := newRW(in)
	p, _ := Stream(r, w)
	out := make([]byte, 6)
	f := func() ([]byte, error) {
		n, err := p.Read(out)
		return out[:n], err
	}
	testRead(t, f, sizes, datas, errors)
}

func TestStreamReadOne(t *testing.T) {
	in := bs("12345\n\n\x00\r")
	sizes := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 0}
	datas := []bs{
		bs("1"), bs("2"), bs("3"), bs("4"), bs("5"),
		bs("\n"), bs("\n"), bs{0}, bs("\r"), bs{},
	}
	errors := []error{nil, nil, nil, nil, nil, nil, nil, nil, nil, io.EOF}

	r, w := newRW(in)
	p, _ := Stream(r, w)
	f := func() ([]byte, error) { return p.ReadOne() }
	testRead(t, f, sizes, datas, errors)
}
