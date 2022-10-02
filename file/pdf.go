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
}

func readSection(r pdfio.Reader, xrefpos int) (*Section, error) {

}

func lastXref(r pdfio.Reader) (*Section, error) {
	r.Seek(-5, io.SeekEnd)
	s, err := r.ReadString(5)
	if err != nil {
		return nil, err
	}

	if s != "%%EOF" {
		return nil, &ErrBadEOF{s}
	}

}

type ErrBadEOF struct {
	is string
}

func (err *ErrBadEOF) Error() string {
	return "bad eof marker: \"" + err.is + "\""
}
