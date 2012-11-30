package eport

import (
	"io"
	"testing"
)

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
