package object

import "github.com/scrouthtv/go-pdf/file"

// Stream is a (virtually unlimited) sequence of bytes.
type Stream struct {
	Blob []byte
}

func ReadStream(r file.Reader) (*Stream, error) {
	panic("not impl")
}
