package object

import (
	"fmt"

	"github.com/scrouthtv/go-pdf/file"
)

type Object interface {
	fmt.Stringer

	Write(file.Writer) error
}

// ReadDirectObject reads anything but an indirect object or a stream.
func ReadDirectObject(r file.Reader) (Object, error) {
	r1, err := r.PeekRune()
	if err != nil {
		return nil, err
	}

	if isNumericCharacter(r1) {
		return ReadNumeric(r)
	}

	switch r1 {
	case 't', 'f':
		return ReadBool(r)
	case '(':
		return ReadString(r)
	case '/':
		return ReadName(r)
	case '[':
		return ReadArray(r)
	case 'n':
		return ReadNull(r)
	case '<':
		// determine whether we have to read a dict or a hex string
		r2, err := r.PeekString(2)
		if err != nil {
			return nil, err
		}

		if r2 == "<<" {
			return ReadDict(r)
		} else {
			return ReadString(r)
		}
	default:
		return nil, &UnexpectedObjectError{r.Position(), "direct object", string(r1)}
	}
}

// ReadArrayMember reads any object that is permissible inside an array.
// Currently this are both direct and indirect objects with the exception of
// streams.
// These are (currently, as of PDF2.0) the same objects as are allowed as
// dictionary value.
func ReadArrayMember(r file.Reader) (Object, error) {
	r1, err := r.PeekRune()
	if err != nil {
		return nil, err
	}

	if isNumericCharacter(r1) {
		if HasIndirect(r) {
			return ReadIndirect(r)
		} else {
			return ReadNumeric(r)
		}
	}

	return ReadDirectObject(r)
}

type UnexpectedObjectError struct {
	Position int
	Expected string
	Got      string
}

func (e *UnexpectedObjectError) Error() string {
	return fmt.Sprintf("unexpected object identifier \"%s\" at %d, expected %s",
		e.Got, e.Position, e.Expected)
}
