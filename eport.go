// Package eport provides a simple API to write Erlang ports in Go.
//
// For more information on Erlang ports see http://www.erlang.org/doc/tutorial/c_port.html.
package eport

import (
	"errors"
)

type Port interface {
	// Read reads up to len(p) bytes into p.
	// It returns the number of bytes read (0 <= n <= len(p)) and any
	// error encountered.
	//
	// In case of line (or packet) based port, receiving a line (packet) of size
	// bigger than len(p) will skip the line (packet) and return n=0 and err=ErrTooBig.
	// It is up to the caller to choose whether to consider this as a fatal error or to continue
	// reading from the port.
	Read(p []byte) (n int, err error)
	// ReadOne reads one packet (or line, if it's a line-based port).
	ReadOne() (data []byte, err error)
	// Write writes len(data) bytes from data to the port.
	// It returns the number of bytes written from data (0 <= n <= len(data))
	// and any error encountered that caused the write to stop early.
	// Write must return a non-nil error if it returns n < len(data).
	Write(data []byte) (n int, err error)
}

var (
	ErrBadSizeLen   = errors.New("eport: bad 'packet size' length")
	ErrSizeOverflow = errors.New("eport: packet size overflow")
	ErrTooBig       = errors.New("eport: packet does not fit the buffer")
)
