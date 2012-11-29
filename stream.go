package eport

import (
	"io"
)

type streamPort struct {
	r io.Reader
	w io.Writer
}

// Stream returns a new stream-based port.
// Note that ReadOne() will only read one byte on each call.
func Stream(r io.Reader, w io.Writer) (Port, error) {
	return &streamPort{r, w}, nil
}

func (p *streamPort) Read(out []byte) (int, error) {
	return p.r.Read(out)
}

func (p *streamPort) ReadOne() ([]byte, error) {
	b := []byte{0}
	_, err := p.r.Read(b)
	return b, err
}

func (p *streamPort) Write(data []byte) (int, error) {
	return p.w.Write(data)
}
