package shared

import (
	"fmt"

	"github.com/scrouthtv/go-pdf/file"
)

type Body interface {
	Resolve(ID) Object
}

type ID interface {
	Equal(other ID) bool
}

type Object interface {
	fmt.Stringer

	Write(file.Writer) error
}
