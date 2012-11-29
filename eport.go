// Package eport provides a simple API to write Erlang ports in Go.
//
// For more information on Erlang ports see http://www.erlang.org/doc/tutorial/c_port.html.
package eport

import (
	"errors"
	"io"
)

// Port is an abstraction over different types of ports.
//
// While using Read() on a line (or packet) based port it may happen that the size
// of the line (packet) is bigger than len(p). In this case eport skips the line (packet)
// and return n=0 and err=ErrTooBig.
// It is up to the caller to choose whether to consider this as a fatal error
// or to continue reading.
type Port interface {
	io.Reader
	io.Writer
	// ReadOne reads either one packet, one line (ending with '\n') or a byte
	// from a packet, line or stream-based port accordingly.
	ReadOne() (data []byte, err error)
}

var (
	ErrBadSizeLen   = errors.New("eport: bad 'packet size' length")
	ErrSizeOverflow = errors.New("eport: packet size overflows integer type")
	ErrTooBig       = errors.New("eport: packet does not fit the buffer")
)
