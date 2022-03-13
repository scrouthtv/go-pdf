package file

import "io"

type Reader interface {
	io.Reader
	io.RuneReader
	io.Seeker
	Position() int

	// ReadString should read a string with len bytes.
	// The specified length is the amount of bytes, *not*
	// the amount of runes.
	ReadString(len int) (string, error)
	PeekRune() (r rune, err error)
	PeekString(len int) (string, error)
}
