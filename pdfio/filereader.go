package pdfio

import (
	"io"
	"os"
	"unicode/utf8"
)

type FileReader struct {
	*os.File
}

func NewReader(f *os.File) *FileReader {
	return &FileReader{f}
}

func (r *FileReader) Peek(buf []byte) (int, error) {
	n, err := r.Read(buf)
	if err != nil {
		return n, err
	}

	_, err = r.Seek(-int64(n), io.SeekCurrent)
	return n, err
}

func (r *FileReader) ReadRune() (rune, int, error) {
	var err error
	var buf []byte
	var i int

	for i = 1; i <= 4; i++ {
		buf = make([]byte, i)
		_, err = r.Peek(buf)
		if err != nil {
			return 0, 0, err
		}

		// TODO only utf8
		// TODO what if we start inside a rune, e.g. RuneStart(buf[0]) == false
		if utf8.FullRune(buf) {
			_, err = r.Seek(int64(i), io.SeekCurrent)
			if err != nil {
				return 0, 0, err
			}

			rn, n := utf8.DecodeRune(buf)
			return rn, n, nil
		}
	}

	panic("this is unexpected")
	return 0, 0, err
}

func (r *FileReader) PeekRune() (rune, error) {
	rn, n, err := r.ReadRune()
	if err != nil {
		return rn, err
	}

	_, err = r.Seek(-int64(n), io.SeekCurrent)
	return rn, err
}

func (r *FileReader) ReadString(length int) (string, error) {
	buf := make([]byte, length)
	_, err := r.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (r *FileReader) PeekString(length int) (string, error) {
	s, err := r.ReadString(length)
	if err != nil {
		return s, err
	}

	_, err = r.Seek(-int64(length), io.SeekCurrent)

	return s, err
}

func (r *FileReader) Position() int {
	pos, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	return int(pos)
}
