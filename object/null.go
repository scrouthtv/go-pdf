package object

import (
	"fmt"

	"github.com/scrouthtv/go-pdf/file"
)

type Null struct{}

var TheNull = Null{}

func ReadNull(r file.Reader) (Null, error) {
	a, err := r.ReadString(4)
	if err != nil {
		return TheNull, err // TODO pack error
	}

	if a == "null" {
		return TheNull, nil
	} else {
		return TheNull, &ErrBadNull{r.Position(), a}
	}
	// fixed length representation,
	// no need to unread anything
}

func (n Null) Write(w file.Writer) error {
	_, err := w.WriteString("null")
	return err
}

func (n Null) String() string {
	return "null"
}

type ErrBadNull struct {
	Position int
	Text     string
}

func (e *ErrBadNull) Error() string {
	return fmt.Sprintf("expected null, got \"%s\" at %d", e.Text, e.Position)
}
