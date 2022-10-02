package pdfio_test

import (
	"go-pdf/pdfio"
	"go-pdf/testutil"
	"io"
	"testing"
)

func TestPrevLine(t *testing.T) {
	f := `startxref
12768
%%EOF`

	pdf := testutil.NewPdf(f)
	pdf.Seek(-2, io.SeekEnd)

	s, err := pdfio.ReadPreviousLine(pdf)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if s != "12768" {
		t.Errorf("bad line read: expected \"%s\", got \"%s\"", "12768", s)
	}
}
