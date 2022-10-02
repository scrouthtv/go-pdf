package pdfio

import (
	"io"
)

type Reader interface {
	io.Reader
	io.RuneReader
	io.Seeker
	Position() int

	// ReadString should read a string with len bytes.
	// The specified length is the amount of bytes, *not*
	// the amount of runes.
	ReadString(length int) (string, error)

	// Peek is like Read(), except that it unreads the contents
	// that were read
	Peek([]byte) (int, error)
	PeekRune() (r rune, err error)
	PeekString(length int) (string, error)

	Advance(amount int) error
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

func ReadPreviousLine(r Reader) (string, error) {
	unreadToEOL(r)
	end := r.Position()

	unreadToEOL(r)
	DiscardEOL(r)
	start := r.Position()

	s, err := r.ReadString(end - start)
	return s, err
}

func unreadToEOL(r Reader) error {
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
