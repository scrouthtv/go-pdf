package file

import (
	"go-pdf/body"
	"go-pdf/header"
	"go-pdf/pdfio"
	"io"
	"strconv"
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

	return &f, nil
}

func readSection(r pdfio.Reader, xrefpos int) (*Section, error) {
	panic("todo")
	return nil, nil
}

// lastXref looks up the position of the last cross-reference table.
func lastXref(r pdfio.Reader) (int, error) {
	r.Seek(0, io.SeekEnd)
	println(r.Position())
	pdfio.UnreadToEOL(r) // TODO: is there *always* an EOL after the EOF marker????

	println(r.Position())

	r.Seek(-5, io.SeekCurrent)
	s, err := r.ReadString(5)
	if err != nil {
		return -1, err
	}

	if s != "%%EOF" {
		return -1, &ErrBadEOF{s}
	}

	xrefpos, err := pdfio.ReadPreviousLine(r)
	if err != nil {
		return -1, err
	}

	label, err := pdfio.ReadPreviousLine(r)
	if err != nil {
		return -1, err
	}

	if label != "startxref" {
		return -1, &ErrBadTrailer{label}
	}

	pos, err := strconv.Atoi(xrefpos)
	if err != nil {
		return -1, err
	}

	return pos, nil
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
