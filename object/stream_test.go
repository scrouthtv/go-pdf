package object_test

import (
	"testing"

	"github.com/scrouthtv/go-pdf/object"
)

func TestStream(t *testing.T) {
	is := `5 0 obj<<
	/THisISAKEy (THisISAVAlue)
	/Length 11
	>>stream
Lorem ipsum
endstreamendobj
	`
	pdf := NewPdf(is)
	i, err := object.ReadIndirect(pdf, nil)
	if err != nil {
		panic(err)
	}
	if i.String() != "indirect(5/0):Stream(60-71:Lo:um)" {
		t.Error(i.String())
	}
}

/*func TestStreamEx(t *testing.T) {
	is := `[7 0 obj
	<</Length 8 0 R>> %An indirect reference to object 8
stream
	BT
		/F1 12 Tf
		72 712 Td
		(A stream with an indirect length) Tj
	ET
endstream
endobj
8 0 obj
	77 %The length of the preceding stream
endobj]`
	pdf := NewPdf(is)
	i, err := object.ReadArray(pdf, nil) // TODO body
	if err != nil {
		panic(err)
	}
	println(i.String())
}
*/
