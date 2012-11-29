// Package eport provides a simple API to write Erlang ports in Go.
//
// For more information on Erlang ports see http://www.erlang.org/doc/tutorial/c_port.html.
package eport

import (
	"errors"
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
