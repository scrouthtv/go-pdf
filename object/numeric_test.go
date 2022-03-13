package object_test

import (
	"testing"

	"github.com/scrouthtv/go-pdf/object"
)

func TestReadInteger(t *testing.T) {
	s := "123 43445 +17 -98 0"
	should := []int{123, 43445, 17, -98, 0}
	pdf := NewPdf(s)

	for i := 0; i < len(should); i++ {
		is, err := object.ReadInteger(pdf)
		if err != nil {
			t.Errorf("error reading %d-nth integer: %s", i, err)
			t.FailNow()
		}

		if int(is) != should[i] {
			t.Errorf("got wrong int at position %d, expected %d, got %d",
				i, should[i], int(is))
		}

		// discard the space:
		pdf.ReadRune()
		// FIXME what about extra whitespaces in pdf?
		// do we have to ignore those as well?
	}
}
