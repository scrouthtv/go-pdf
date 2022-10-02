package object_test

import (
	"testing"

	"go-pdf/object"
)

func TestBasicDict(t *testing.T) {
	in := `<</Type /Example
	/Subtype /DictionaryExample
	/Version 0.01
	/Integer#20Item 12
	/StringItem (a string)/StringItem#7b <3A20>
	/Subdictionary <<
		/Item1 0.4
		/Item2 true
		/LastItem (not !)
		/VeryLastItem (OK)
	>>
>>`

	pdf := NewPdf(in)

	d, err := object.ReadDict(pdf, nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(d)
}
