package eport

import (
	"bufio"
	"io"
)

type linePort struct {
	r *bufio.Reader
	w io.Writer
}

// Line returns a new line-based port.
func Line(r io.Reader, w io.Writer) (Port, error) {
	return &linePort{
		bufio.NewReader(r),
		w,
	}, nil
}

func (p *linePort) Read(out []byte) (int, error) {
	line, err := p.r.ReadSlice('\n')
	size := len(line)

	if size > len(out) {
		return 0, ErrTooBig
	}

	copy(out, line)

	return size, err
}

func (p *linePort) ReadOne() ([]byte, error) {
	return p.r.ReadBytes('\n')
}

func (p *linePort) Write(data []byte) (int, error) {
	return p.w.Write(data)
}
