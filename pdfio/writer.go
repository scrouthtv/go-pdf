package file

import "io"

type Writer interface {
	io.Writer
	io.StringWriter
	WriteRune(rune) error
}
