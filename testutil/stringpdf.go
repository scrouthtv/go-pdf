package testutil

import (
	"io"
	"strings"
	"testing"

	"go-pdf/object"
	"go-pdf/pdfio"
)

type StringPDF struct {
	strings.Reader
}

func NewPdf(s string) *StringPDF {
	return &StringPDF{*strings.NewReader(s)}
}

func (s *StringPDF) Position() int {
	return int(s.Size()) - s.Len()
}

func (s *StringPDF) Peek(buf []byte) (int, error) {
	n, err := s.Read(buf)
	if err != nil {
		return 0, err
	}

	_, err = s.Seek(-int64(n), io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	return n, err
}

func (s *StringPDF) PeekRune() (rune, error) {
	read, n, err := s.ReadRune()
	if err != nil {
		return 0, err
	}

	for n > 0 {
		err = s.UnreadByte()
		if err != nil {
			return 0, err
		}
		n--
	}

	return read, nil
}

func (s *StringPDF) PeekString(length int) (string, error) {
	buf := make([]byte, length)
	_, err := s.Read(buf)
	if err != nil {
		return "", err
	}

	// unread string:
	for length > 0 {
		err = s.UnreadByte()
		if err != nil {
			return "", err
		}
		length--
	}

	return string(buf), nil
}

func (s *StringPDF) ReadString(length int) (string, error) {
	buf := make([]byte, length)

	_, err := s.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (s *StringPDF) Advance(amount int) error {
	_, err := s.Seek(int64(amount), io.SeekCurrent)
	if err != nil {
		return err
	}
	return nil
}

func TestDiscardWhitespace(t *testing.T) {
	s := "2  15 \t a\r\nb"
	pdf := NewPdf(s)

	expectAfterSeek(t, pdf, 0)

	discardRune(t, pdf)

	expectAfterSeek(t, pdf, 3)

	discardRune(t, pdf)
	discardRune(t, pdf)

	expectAfterSeek(t, pdf, 8)

	discardRune(t, pdf)

	expectAfterSeek(t, pdf, 11)
}

func expectAfterSeek(t *testing.T, pdf pdfio.Reader, pos int) {
	t.Helper()

	err := object.DiscardWhitespace(pdf)
	if err != nil {
		t.Error(err)
	}

	if pdf.Position() != pos {
		t.Errorf("expected position %d, got %d", pos, pdf.Position())
	}
}

func discardRune(t *testing.T, pdf pdfio.Reader) {
	t.Helper()

	_, _, err := pdf.ReadRune()
	if err != nil {
		t.Error(err)
	}
}
