package port

import (
	"io"
	"testing"
)

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
