// Package eport provides a simple API to write Erlang ports in Go.
//
// For more information on Erlang ports see http://www.erlang.org/doc/tutorial/c_port.html.
package eport

import (
	"bufio"
	. "encoding/binary"
	"errors"
	"io"
)

type Port interface {
	// ReadOne reads one packet (or line, if it's a line-based port).
	ReadOne() (data []byte, err error)
	// Write writes len(data) bytes from data to the port.
	// It returns the number of bytes written from data (0 <= n <= len(data))
	// and any error encountered that caused the write to stop early.
	// Write must return a non-nil error if it returns n < len(data).
	Write(data []byte) (n int, err error)
}

var (
	ErrSizeOverflow = errors.New("eport: packet size overflow")
	ErrBadSizeLen   = errors.New("eport: bad 'packet size' length")
)

// Line returns a new line-based port.
func Line(r io.Reader, w io.Writer) (Port, error) {
	return &linePort{
		bufio.NewReader(r),
		w,
	}, nil
}

// Packet returns a new packet-based port.
// Each packet is preceded with its length.
// The size of the length of the packet is either 1, 2 or 4.
func Packet(r io.Reader, w io.Writer, sizeLen int) (Port, error) {
	switch sizeLen {
	case 1, 2, 4:
		return &packetPort{r, w, 1 << (uint(sizeLen) * 8), make([]byte, sizeLen)}, nil
	}

	return nil, ErrBadSizeLen
}

// Stream returns a new stream-based port.
// Note that ReadOne() will only read one byte on each call.
func Stream(r io.Reader, w io.Writer) (Port, error) {
	return &streamPort{r, w}, nil
}

type linePort struct {
	r *bufio.Reader
	w io.Writer
}

func (p *linePort) ReadOne() ([]byte, error) {
	return p.r.ReadBytes('\n')
}

func (p *linePort) Write(data []byte) (int, error) {
	return p.w.Write(data)
}

type packetPort struct {
	r       io.Reader
	w       io.Writer
	max     int
	sizeBuf []byte
}

func (pr *packetPort) ReadOne() ([]byte, error) {
	if _, err := io.ReadFull(pr.r, pr.sizeBuf); err != nil {
		return nil, err
	}

	var size int

	switch len(pr.sizeBuf) {
	case 1:
		size = int(pr.sizeBuf[0])
	case 2:
		size = int(BigEndian.Uint16(pr.sizeBuf))
	case 4:
		size32 := BigEndian.Uint32(pr.sizeBuf)
		size = int(size32)
		if uint32(size) != size32 {
			return nil, ErrSizeOverflow
		}
	}

	if size == 0 {
		return []byte{}, nil
	}

	data := make([]byte, size)
	if _, err := io.ReadFull(pr.r, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (pr *packetPort) Write(data []byte) (int, error) {
	size := len(data)
	if size > pr.max {
		return 0, ErrSizeOverflow
	}

	switch len(pr.sizeBuf) {
	case 1:
		pr.sizeBuf[0] = uint8(size)
	case 2:
		BigEndian.PutUint16(pr.sizeBuf, uint16(size))
	case 4:
		BigEndian.PutUint32(pr.sizeBuf, uint32(size))
	}

	if n, err := pr.w.Write(pr.sizeBuf); err != nil {
		return n, err
	}

	return pr.w.Write(data)
}

type streamPort struct {
	r io.Reader
	w io.Writer
}

func (p *streamPort) ReadOne() ([]byte, error) {
	b := []byte{0}
	_, err := p.r.Read(b)
	return b, err
}

func (p *streamPort) Write(data []byte) (int, error) {
	return p.w.Write(data)
}
