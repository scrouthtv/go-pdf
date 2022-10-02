package header

import "go-pdf/pdfio"

type Header struct {
	Ver            Version
	ContainsBinary bool
}

func ReadHeader(r pdfio.Reader) (*Header, error) {
	var h Header

	v, err := ReadVersion(r)
	if err != nil {
		return nil, err
	}

	h.Ver = v

	// Check for binary data:
	is, comment, err := pdfio.ReadComment(r)
	if err != nil {
		return nil, err
	}

	if !is {
		h.ContainsBinary = false
		return &h, nil
	}

	h.ContainsBinary = isBinaryMarker(comment)

	return &h, nil
}

func isBinaryMarker(comment string) bool {
	bchars := 0

	for _, r := range comment {
		if r >= 128 {
			bchars++
		}
	}

	return bchars >= 4
}
