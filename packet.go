package eport

import (
	. "encoding/binary"
	"io"
)

type packetPort struct {
	r       io.Reader
	w       io.Writer
	max     int
	sizeBuf []byte
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
