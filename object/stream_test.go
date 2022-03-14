package object_test

import (
	"testing"

	"github.com/scrouthtv/go-pdf/body"
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
		t.Error(err)
	}

	if i.String() != "indirect(5/0):Stream(60-71)" {
		t.Error("got bad stream data:", i.String())
	}
}

// TODO comments

func TestStreamEx(t *testing.T) {
	is := `[8 0 obj
	77
endobj
<</Length 8 0 R>>
stream
	BT
		/F1 12 Tf
		72 712 Td
		(A stream with an indirect length) Tj
	ET
endstream
endobj
]
`

	pdf := NewPdf(is)
	i, err := object.ReadArray(pdf, body.NewBody()) // TODO body
	if err != nil {
		t.Error(err)
	}

	println(i.String())
}
