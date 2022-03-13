package object

import (
	"fmt"

	"github.com/scrouthtv/go-pdf/file"
)

type Bool bool

var BadBool = Bool(false)

func ReadBool(r file.Reader) (Bool, error) {
	a, err := r.ReadString(4)
	if err != nil {
		return BadBool, err // TODO pack error
	}

	switch a {
	case "true":
		return Bool(true), nil
	case "fals":
		read, _, err := r.ReadRune()
		if err != nil {
			return BadBool, err // TODO pack error
		}

		if read == 'e' {
			return Bool(false), nil
		} else {
			return BadBool, &ErrBadBool{r.Position(), a + string(read)}
		}
	default:
		return BadBool, &ErrBadBool{r.Position(), a}
	}

	// fixed length representations,
	// don't need to unread anything
}

func (b Bool) Write(w file.Writer) error {
	var err error

	switch b {
	case true:
		_, err = w.WriteString("true")
	case false:
		_, err = w.WriteString("false")
	}

	return err
}

func (b Bool) String() string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

type ErrBadBool struct {
	Position int
	Text     string
}

func (e *ErrBadBool) Error() string {
	return fmt.Sprintf("expected true/false, got \"%s\" at %d", e.Text, e.Position)
}
