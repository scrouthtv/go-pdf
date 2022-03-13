package object_test

import (
	"strings"
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
