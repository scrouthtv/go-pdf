package object

import (
	"fmt"

	"github.com/scrouthtv/go-pdf/file"
)

type Object interface {
	fmt.Stringer

	Write(file.Writer) error
}

func ReadObject(r file.Reader) (Object, error) {
	// TODO
	panic("not impl")
}
