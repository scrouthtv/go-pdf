package object

import (
	"fmt"
	"strings"

	"github.com/scrouthtv/go-pdf/file"
	"github.com/scrouthtv/go-pdf/shared"
)

// Array is a collection of objects.
// PDF arrays may contain objects of any (and mixed) type.
// This also allows for nested arrays.
type Array struct {
	Elems []shared.Object
}

var BadArray = &Array{}

func ReadArray(r file.Reader, b shared.Body) (*Array, error) {
	id, _, err := r.ReadRune()
	if err != nil {
		return nil, &MissingArrayTokenError{r.Position(), err}
	}

	if id != '[' {
		return nil, &BadArrayStartError{r.Position(), id}
	}

	id, err = r.PeekRune()
	if err != nil {
		return nil, &MissingArrayTokenError{r.Position(), err}
	}

	arr := Array{}

	for id != ']' {
		obj, err := ReadArrayMember(r, b)
		if err != nil {
			return BadArray, &BadArrayMemberError{err}
		}

		println("got a member:", obj.String())

		arr.Elems = append(arr.Elems, obj)

		DiscardWhitespace(r)

		id, err = r.PeekRune()
		if err != nil {
			return BadArray, &MissingArrayTokenError{r.Position(), err}
		}
	}

	// end is detected by the closing character,
	// so we don't need to unread the final character.

	return &arr, nil
}

func (a *Array) Write(w file.Writer) error {
	err := w.WriteRune('[')
	if err != nil {
		return err
	}

	for i, e := range a.Elems {
		if i > 0 {
			err = w.WriteRune(' ')
			if err != nil {
				return err
			}
		}

		err = e.Write(w)
		if err != nil {
			return err
		}
	}

	return w.WriteRune(']')
}

func (a *Array) String() string {
	var out strings.Builder

	out.WriteString("[")

	for i, obj := range a.Elems {
		if i > 0 {
			out.WriteString(", ")
		}

		out.WriteString(obj.String())
	}

	return out.String()
}

type BadArrayStartError struct {
	Pos   int
	Start rune
}

func (e *BadArrayStartError) Error() string {
	return fmt.Sprintf("expected array at pos %d, got %q instead", e.Pos, e.Start)
}

type RunawayArrayMemberError struct {
	Pos  int
	What rune
}

func (e *RunawayArrayMemberError) Error() string {
	return fmt.Sprintf("runaway array member at pos %d, got %q", e.Pos, e.What)
}

type BadArrayMemberError struct {
	Err error
}

func (e *BadArrayMemberError) Error() string {
	return "encountered an error reading array member: " + e.Err.Error()
}

func (e *BadArrayMemberError) Unwrap() error {
	return e.Err
}

type MissingArrayTokenError struct {
	Position int
	Err      error
}

func (e *MissingArrayTokenError) Error() string {
	return fmt.Sprintf("error reading array token at %d: %s", e.Position, e.Err.Error())
}
