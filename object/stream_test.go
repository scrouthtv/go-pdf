package object_test

import (
	"testing"

	"go-pdf/body"
	"go-pdf/object"
	"go-pdf/pdfio"
	"go-pdf/testutil"
)

func TestStream(t *testing.T) {
	is := `5 0 obj<<
	/THisISAKEy (THisISAVAlue)
	/Length 11
	>>stream
Lorem ipsum
endstreamendobj
	`
	pdf := testutil.NewPdf(is)
	i, err := object.ReadIndirect(pdf, nil)
	if err != nil {
		t.Error(err)
	}

	if i.String() != "indirect(5/0):Stream(60-71)" {
		t.Error("got bad stream data:", i.String())
	}
}

// TODO comments

func TestStreamBlind(t *testing.T) {
	in := `stream
asdf qwertz
endstream`

	pdf := testutil.NewPdf(in)

	is, err := object.ReadStream(pdf, object.NewDict())
	if err != nil {
		t.Error(err)
	}

	t.Log(is)
}

func TestStreamEx(t *testing.T) {
	in := `8 0 obj
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
`
	pdf := testutil.NewPdf(in)
	b := body.NewBody()

	is, err := object.ReadIndirect(pdf, b)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	pdfio.DiscardWhitespace(pdf)
	if pdf.Position() != 19 {
		t.Errorf("Wrong position after first indirect, expected 19, got %d", pdf.Position())
	}

	b.Obj = append(b.Obj, is.(*object.IndirectVal))

	isd, err := object.ReadIndirect(pdf, b) // TODO body
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isd)
}
