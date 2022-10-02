package object

import (
	"fmt"
	"io"

	"go-pdf/pdfio"
)

// Name is an (internal) name object.
// Names are atomic (have no structure) and
// are unique to the document.
// PDF name objects are compared byte-by-byte,
// even if the UTF8 representation is equal.
type Name struct {
	Name string
}

func ReadName(r pdfio.Reader) (*Name, error) {
	read, _, err := r.ReadRune()
	if err != nil {
		return nil, err // TODO pack error
	}

	if read != '/' {
		return nil, &BadNameStartError{r.Position(), read}
	}

	out := Name{""}

	for {
		read, _, err = r.ReadRune()
		if err != nil {
			return nil, err // TODO pack error
		}

		if read == '#' {
			read, err = readHexCharacter(r)
			if err != nil {
				return nil, err // TODO pack error
			}

			out.Name += string(read)
		} else if isTokenDelimiter(read) {
			// end is detected by an invalid character,
			// so we have to unread the invalid character:
			r.Seek(-1, io.SeekCurrent)
			return &out, nil
		} else {
			out.Name += string(read)
		}
	}
}

func (n *Name) Write(w pdfio.Writer) error {
	err := w.WriteRune('/')
	if err != nil {
		return err
	}

	// range automatically breaks out entire runes, not bytes.
	for _, r := range n.Name {
		if isRegularCharacter(r) {
			err = w.WriteRune(r)
			if err != nil {
				return err
			}
		} else {
			panic("not impl")
		}
	}

	return nil
}

func (n *Name) String() string {
	return "name(" + n.Name + ")"
}

type BadNameStartError struct {
	Pos   int
	Start rune
}

func (e *BadNameStartError) Error() string {
	return fmt.Sprintf("expected name at pos %d, got %q instead", e.Pos, e.Start)
}
