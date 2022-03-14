package object_test

import (
	"testing"

	"github.com/scrouthtv/go-pdf/object"
)

func TestBasicDict(t *testing.T) {
	in := `<</Type /Example
	/Subtype /DictionaryExample
	/Version 0.01
	/IntegerItem 12
	/StringItem (a string)
	/Subdictionary <<
		/Item1 0.4
		/Item2 true
		/LastItem (not !)
		/VeryLastItem (OK)
	>>
>>`

	pdf := NewPdf(in)

	d, err := object.ReadDict(pdf)
	if err != nil {
		t.Error(err)
	}

	t.Log(d)
}
