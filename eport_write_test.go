package eport

import (
	"testing"
)

func TestLineWrite(t *testing.T) {
	final := bs("1234567890123\n\nфыва\r\n\tолдж\n1")
	outs := []bs{bs("123456789"), bs("0123\n"), bs("\nфыва"), bs("\r\n\tолдж\n1")}
	sizes := []int{9, 5, 9, 13}

	r, w := newRW(bs{})
	p, _ := Line(r, w)
	testWrite(t, p, w, sizes, outs, final)
}

func TestPacket1Write(t *testing.T) {
	final := bs("\x09123456789\x050123\n\x09\nфыва\x0d\r\n\tолдж\n1")
	outs := []bs{bs("123456789"), bs("0123\n"), bs("\nфыва"), bs("\r\n\tолдж\n1")}
	sizes := []int{9, 5, 9, 13}

	r, w := newRW(bs{})
	p, _ := Packet(r, w, 1)
	testWrite(t, p, w, sizes, outs, final)

	exp := 255
	n, err := p.Write(make([]byte, exp))
	if n != exp {
		t.Errorf("expected %d, got %d", exp, n)
	}
	if err != nil {
		t.Errorf("err is %#v, expected nil", err)
	}
}

func TestPacket2Write(t *testing.T) {
	final := bs("\x00\x09123456789\x00\x050123\n\x00\x09\nфыва\x00\x0d\r\n\tолдж\n1")
	outs := []bs{bs("123456789"), bs("0123\n"), bs("\nфыва"), bs("\r\n\tолдж\n1")}
	sizes := []int{9, 5, 9, 13}

	r, w := newRW(bs{})
	p, _ := Packet(r, w, 2)
	testWrite(t, p, w, sizes, outs, final)

	exp := 65535
	n, err := p.Write(make([]byte, exp))
	if n != exp {
		t.Errorf("expected %d, got %d", exp, n)
	}
	if err != nil {
		t.Errorf("err is %#v, expected nil", err)
	}
}

func TestPacket4Write(t *testing.T) {
	final := bs("\x00\x00\x00\x09123456789\x00\x00\x00\x050123\n\x00\x00\x00\x09\nфыва\x00\x00\x00\x0d\r\n\tолдж\n1")
	outs := []bs{bs("123456789"), bs("0123\n"), bs("\nфыва"), bs("\r\n\tолдж\n1")}
	sizes := []int{9, 5, 9, 13}

	r, w := newRW(bs{})
	p, _ := Packet(r, w, 4)
	testWrite(t, p, w, sizes, outs, final)
}

func TestPacket1SizeOverflow(t *testing.T) {
	r, w := newRW(bs{})
	p, _ := Packet(r, w, 1)

	n, err := p.Write(make([]byte, 256))
	if n != 0 {
		t.Errorf("expected %d, got %d", 0, n)
	}
	if err != ErrSizeOverflow {
		t.Errorf("expected %#v, got %#v", ErrSizeOverflow, err)
	}
}

func TestPacket2SizeOverflow(t *testing.T) {
	r, w := newRW(bs{})
	p, _ := Packet(r, w, 2)

	n, err := p.Write(make([]byte, 65536))
	if n != 0 {
		t.Errorf("expected %d, got %d", 0, n)
	}
	if err != ErrSizeOverflow {
		t.Errorf("expected %#v, got %#v", ErrSizeOverflow, err)
	}
}

func TestStreamWrite(t *testing.T) {
	final := bs("1234567890123\n\nфыва\r\n\tолдж\n1")
	outs := []bs{bs("123456789"), bs("0123\n"), bs("\nфыва"), bs("\r\n\tолдж\n1")}
	sizes := []int{9, 5, 9, 13}

	r, w := newRW(bs{})
	p, _ := Stream(r, w)
	testWrite(t, p, w, sizes, outs, final)
}
