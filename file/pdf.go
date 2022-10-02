package file

import (
	"go-pdf/header"
)

type File struct {
	hdr      header.Header
	sections []Section
}

type Section struct {
	//body body.Body
	xref header.XRef
	trlr header.Trailer
}
