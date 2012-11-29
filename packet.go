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

func (p *packetPort) Read(out []byte) (int, error) {
	size, err := p.readSize()
	if err != nil {
		return 0, err
	} else if size == 0 {
		return 0, nil
	}

	if size > len(out) {
		// skip the packet, so too big packets can be ignored
		for size > 0 && (err == nil || err == io.EOF) {
			var n int
			n, err = p.r.Read(out)
			size -= n
		}

		return 0, ErrTooBig
	}

	return size, nil
}

func (p *packetPort) ReadOne() ([]byte, error) {
	size, err := p.readSize()
	if err != nil {
		return nil, err
	} else if size == 0 {
		return []byte{}, nil
	}

	data := make([]byte, size)
	_, err := io.ReadFull(p.r, data)

	return data, err
}

func (p *packetPort) Write(data []byte) (int, error) {
	size := len(data)
	if size > p.max {
		return 0, ErrSizeOverflow
	}

	switch len(p.sizeBuf) {
	case 1:
		p.sizeBuf[0] = uint8(size)
	case 2:
		BigEndian.PutUint16(p.sizeBuf, uint16(size))
	case 4:
		BigEndian.PutUint32(p.sizeBuf, uint32(size))
	}

	if n, err := p.w.Write(p.sizeBuf); err != nil {
		return n, err
	}

	return p.w.Write(data)
}

func (p *packetPort) readSize() (int, error) {
	if _, err := io.ReadFull(p.r, p.sizeBuf); err != nil {
		return 0, err
	}

	var size int

	switch len(p.sizeBuf) {
	case 1:
		size = int(p.sizeBuf[0])
	case 2:
		size = int(BigEndian.Uint16(p.sizeBuf))
	case 4:
		size32 := BigEndian.Uint32(p.sizeBuf)
		size = int(size32)
		if uint32(size) != size32 {
			return 0, ErrSizeOverflow
		}
	}

	return size, nil
}
