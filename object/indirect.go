package object

import (
	"fmt"
	"io"

	"github.com/scrouthtv/go-pdf/file"
)

type ObjID struct {
	// ID is the object number. Object numbers can be arbitrarily assigned.
	// However, they should be unique.
	ID Integer

	// Gen is the object's generation number. It is non-negative and
	// defaults to 0. Other generation numbers only appear in later file
	// updates.
	Gen Integer
}

func (o *ObjID) String() string {
	return fmt.Sprintf("%d/%d", o.ID, o.Gen)
}

type Indirect struct {
	ID    ObjID
	Value Object
}

// HasIndirect checks whether the next token is an indirect object.
// If an error occurs during reading, false is returned.
func HasIndirect(r file.Reader) bool {
	defer r.Seek(int64(r.Position()), io.SeekStart)

	spc := 0

	for spc < 2 {
		read, _, err := r.ReadRune()
		if err != nil {
			return false
		}

		if read == ' ' {
			spc++
		}
	}

	s, err := r.ReadString(3)
	if err != nil {
		return false
	}

	return s == "obj"
}

func ReadIndirect(r file.Reader) (*Indirect, error) {
	i := Indirect{}
	var err error

	i.ID.ID, err = ReadInteger(r)
	if err != nil {
		return nil, err
	}

	DiscardWhitespace(r)

	i.ID.Gen, err = ReadInteger(r)
	if err != nil {
		return nil, err
	}

	DiscardWhitespace(r)

	rs, err := r.ReadString(3)
	if err != nil {
		return nil, err
	}

	if rs != "obj" {
		return nil, &BadIndirectSpecifierError{r.Position(), "obj", rs}
	}

	DiscardWhitespace(r) // whitespace after obj

	// FIXME this can also be an indirect object
	i.Value, err = ReadDirectObject(r)
	if err != nil {
		return nil, err
	}

	DiscardWhitespace(r) // whitespace after actual contents

	rs, err = r.ReadString(6)
	if err != nil {
		return nil, err
	}

	if rs != "endobj" {
		return nil, &BadIndirectSpecifierError{r.Position(), "endobj", rs}
	}

	return &i, nil
}

func (i *Indirect) Write(w file.Writer) error {
	err := i.ID.ID.Write(w)
	if err != nil {
		return err
	}

	err = w.WriteRune(' ')
	if err != nil {
		return err
	}

	err = i.ID.Gen.Write(w)
	if err != nil {
		return err
	}

	_, err = w.WriteString(" obj\r\n")
	if err != nil {
		return err
	}

	err = i.Value.Write(w)
	if err != nil {
		return err
	}

	_, err = w.WriteString("\r\nendobj")

	return err
}

func (i *Indirect) String() string {
	return fmt.Sprintf("indirect(%s):%s", i.ID.String(), i.Value.String())
}

type RunawayIndirectMemberError struct {
	Position int
}

func (e *RunawayIndirectMemberError) Error() string {
	return fmt.Sprintf("runaway indirect member, excepted space at %d", e.Position)
}

type BadIndirectSpecifierError struct {
	Position  int
	Expected  string
	Specifier string
}

func (e *BadIndirectSpecifierError) Error() string {
	return fmt.Sprintf("bad indirect specifier, expected %s, got %s at %d",
		e.Specifier, e.Expected, e.Position)
}
