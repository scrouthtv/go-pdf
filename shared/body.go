package shared

import (
	"fmt"

	"go-pdf/pdfio"
)

type Body interface {
	Resolve(ID) Object
}

type ID interface {
	Equal(other ID) bool
}

type Object interface {
	fmt.Stringer

	Write(pdfio.Writer) error
}
