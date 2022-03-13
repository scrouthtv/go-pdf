package object

import (
	"fmt"
	"strings"

	"github.com/scrouthtv/go-pdf/file"
)

// Array is a collection of objects.
// PDF arrays may contain objects of any (and mixed) type.
// This also allows for nested arrays.
type Array struct {
	Elems []Object
}

var BadArray = &Array{}

func ReadArray(r file.Reader) (*Array, error) {
	id, _, err := r.ReadRune()
	if err != nil {
		return nil, err
	}

	if id != '[' {
		return nil, &ErrBadArrayStart{r.Position(), id}
	}

	id, err = r.PeekRune()
	if err != nil {
		return nil, err
	}

	arr := Array{}

	for id != ']' {
		obj, err := ReadArrayMember(r)
		if err != nil {
			return BadArray, err // TODO pack error
		}

		arr.Elems = append(arr.Elems, obj)

		id, _, err = r.ReadRune()
		if err != nil {
			return BadArray, err // TODO pack error
		}

		if id != ' ' && id != ']' {
			what, _ := r.PeekRune()
			return BadArray, &ErrRunawayArrayMember{r.Position(), what}
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

type ErrBadArrayStart struct {
	Pos   int
	Start rune
}

func (e *ErrBadArrayStart) Error() string {
	return fmt.Sprintf("expected array at pos %d, got %q instead", e.Pos, e.Start)
}

type ErrRunawayArrayMember struct {
	Pos  int
	What rune
}

func (e *ErrRunawayArrayMember) Error() string {
	return fmt.Sprintf("runaway array member at pos %d, got %q", e.Pos, e.What)
}
