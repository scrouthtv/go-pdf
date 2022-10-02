package object

import (
	"io"
	"strconv"

	"go-pdf/pdfio"
	"go-pdf/shared"
)

// Numeric may either be an Integer or a Floating number.
type Numeric interface {
	shared.Object
}

// ReadNumeric may read either an Integer or a Real.
func ReadNumeric(r pdfio.Reader) (Numeric, error) {
	a := ""

	read, _, err := r.ReadRune()
	if err != nil {
		return BadInteger, err
	}

	isfloat := false

	for isNumericCharacter(read) {
		a += string(read)

		if read == '.' {
			isfloat = true
		}

		read, _, err = r.ReadRune()
		if err != nil {
			return nil, err
		}
	}

	// end is detected by an invalid character,
	// so we have to unread the invalid character:
	r.Seek(-1, io.SeekCurrent)

	if isfloat {
		return NewFloatFromString(a)
	} else {
		return NewIntegerFromString(a)
	}
}

type Integer int64

var BadInteger = Integer(0)

func NewInteger(i int64) Integer {
	return Integer(i)
}

func NewIntegerFromString(s string) (Integer, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return BadInteger, err
	}

	return Integer(i), nil
}

func ReadInteger(r pdfio.Reader) (Integer, error) {
	a := ""

	read, _, err := r.ReadRune()
	if err != nil {
		return BadInteger, err // TODO pack error
	}

	for read >= '0' && read <= '9' || read == '+' || read == '-' {
		a += string(read)

		read, _, err = r.ReadRune()
		if err != nil {
			continue
		}
	}

	// end is detected by an invalid character,
	// so we have to unread the invalid character:
	r.Seek(-1, io.SeekCurrent)

	return NewIntegerFromString(a)
}

func (i Integer) Write(w pdfio.Writer) error {
	_, err := w.WriteString(strconv.Itoa(int(i)))
	return err
}

func (i Integer) String() string {
	return "integer(" + strconv.Itoa(int(i)) + ")"
}

type Floating float64

var BadFloating = Floating(0)

func NewFloating(f float64) Floating {
	return Floating(f)
}

func NewFloatFromString(s string) (Floating, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return BadFloating, err
	}

	return Floating(f), nil
}

func (f Floating) Write(w pdfio.Writer) error {
	_, err := w.WriteString(strconv.FormatFloat(float64(f), 'f', -1, 64))
	return err
}

func (f Floating) String() string {
	return "floating(" + strconv.FormatFloat(float64(f), 'f', -1, 64) + ")"
}
