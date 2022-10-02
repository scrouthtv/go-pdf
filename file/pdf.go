package file

import (
	"go-pdf/body"
	"go-pdf/header"
	"go-pdf/pdfio"
	"io"
)

type File struct {
	hdr      header.Header
	sections []Section
}

type Section struct {
	body body.Body
	xref header.XRef
	trlr header.Trailer
}

func ReadFile(r pdfio.Reader) (*File, error) {
	var f File

	h, err := header.ReadHeader(r)
	if err != nil {
		return nil, err
	}

	f.hdr = *h

	// Read last xref table:
	moveToEnd(r)
	pos, err := prevXref(r)
	if err != nil {
		return nil, err
	}

	s, err := readSection(r, pos)

	return &f, nil
}

func readSection(r pdfio.Reader, xrefpos int) (*Section, error) {
	panic("todo")
	return nil, nil
}

func moveToEnd(r pdfio.Reader) error {
	r.Seek(0, io.SeekEnd)
	println(r.Position())
	pdfio.UnreadToEOL(r) // TODO: is there *always* an EOL after the EOF marker????

	println(r.Position())

	r.Seek(-5, io.SeekCurrent)
	s, err := r.ReadString(5)
	if err != nil {
		return err
	}

	if s != "%%EOF" {
		return &ErrBadEOF{s}
	}

	return nil
}

// prevXref looks up the position of the previous cross-reference table.
// This is done by reversing line-by-line over the file until "startxref"
// is found. The following line is returned.
func prevXref(r pdfio.Reader) (int, error) {
	label, err := pdfio.ReadPreviousLine(r)
	if err != nil {
		return -1, err
	}

	if label != "startxref" {
		return -1, &ErrBadTrailer{label}
	}

	pdfio.ReadLine
}

type ErrBadEOF struct {
	is string
}

func (err *ErrBadEOF) Error() string {
	return "bad eof marker: \"" + err.is + "\""
}

type ErrBadTrailer struct {
	StartXrefLabel string
}

func (err *ErrBadTrailer) Error() string {
	return "bad trailer, expected \"startxref\", got \"" + err.StartXrefLabel + "\""
}
