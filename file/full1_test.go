package file

import (
	"go-pdf/pdfio"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	f, err := os.Open("simple.pdf")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	r := pdfio.NewReader(f)

	pos, err := prevXref(r)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if pos != 3808 {
		t.Errorf("bad xref position: expected 3808, got %d", pos)
	} else {
		t.Log("correct xref position")
	}
}
