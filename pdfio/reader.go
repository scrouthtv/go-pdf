package pdfio

import (
	"io"
)

type Reader interface {
	io.Reader
	io.RuneReader
	io.Seeker
	Position() int // TODO this limits the position to 2.1 GB

	// ReadString should read a string with len bytes.
	// The specified length is the amount of bytes, *not*
	// the amount of runes.
	ReadString(length int) (string, error)

	// Peek is like Read(), except that it unreads the contents
	// that were read
	Peek([]byte) (int, error)
	PeekRune() (r rune, err error)
	PeekString(length int) (string, error)
}

// ReadComment attempts to read a comment from the pdf.
// It expects the next rune in the file to be a '%'.
// If it is not, it returns false and an empty string.
// The extra error is returned in case an error occurs
// during io.
func ReadComment(r Reader) (bool, string, error) {
	rn, n, err := r.ReadRune()
	if err != nil {
		return false, "", err
	}

	if rn != '%' {
		_, err = r.Seek(-int64(n), io.SeekCurrent)
		if err != nil {
			return false, "", err
		}

		return false, "", nil
	}

	rn, _, err = r.ReadRune()
	if err != nil {
		return false, "", err
	}

	cmt := ""
	for rn != '\n' {
		cmt = cmt + string(rn)

		rn, _, err = r.ReadRune()
		if err != nil {
			return false, "", err
		}
	}

	return true, cmt, nil
}

// ReadToEOL reads up to and including the next EOL
// and returns the string without the EOL marker.
func ReadToEOL(r Reader) (string, error) {
	var out string
	buf := make([]byte, 1)

	for {
		_, err := r.Read(buf)
		if err != nil {
			return "", err
		}

		if IsEOL(rune(buf[0])) { // FIXME check for CRLF
			return out, nil
		} else {
			out = out + string(buf)
		}
	}
}

// ReadPreviousLine reads the previous line and returns it without
// preceeding or trailing EOL markers.
// After the call, the file is positioned after the EOL of the line.
func ReadPreviousLine(r Reader) (string, error) {
	UnreadToEOL(r)
	end := r.Position()

	UnreadToEOL(r)
	DiscardEOL(r)
	start := r.Position()

	println("line from ", start, "to", end)

	s, err := r.ReadString(end - start)
	if err != nil {
		return s, err
	}

	err = DiscardEOL(r)
	return s, err
}

func UnreadToEOL(r Reader) error {
	var buf []byte = make([]byte, 1)

	for {
		r.Seek(-1, io.SeekCurrent)
		_, err := r.Peek(buf)
		if err != nil {
			return err
		}

		if buf[0] == '\n' {
			// Check if CRLF:
			r.Seek(-1, io.SeekCurrent)
			_, err := r.Read(buf)
			if err != nil {
				return err
			}

			// if CRLF unread to beginning of CRLF:
			if buf[0] == '\r' {
				r.Seek(-1, io.SeekCurrent)
			}

			return nil
		} else if buf[0] == '\r' {
			return nil
		}
	}
}
