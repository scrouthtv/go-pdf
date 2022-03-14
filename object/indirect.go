package object

import (
	"fmt"
	"io"

	"github.com/scrouthtv/go-pdf/file"
	"github.com/scrouthtv/go-pdf/shared"
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

func (o *ObjID) Equal(other shared.ID) bool {
	return false // TODO
}

type Indirect interface {
	shared.Object
	Value() shared.Object
}

type IndirectVal struct {
	ID    ObjID
	value shared.Object
}

type IndirectRef struct {
	ID   ObjID
	body shared.Body
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

	bR, err := r.PeekRune()
	if err != nil {
		return false
	}
	if bR == 'R' {
		r.ReadRune()
		return true
	}

	s, err := r.ReadString(3)
	if err != nil {
		return false
	}

	return s == "obj"
}

func ReadIndirect(r file.Reader, b shared.Body) (Indirect, error) {
	i := ObjID{}
	var err error

	i.ID, err = ReadInteger(r)
	if err != nil {
		return nil, err
	}

	DiscardWhitespace(r)

	i.Gen, err = ReadInteger(r)
	if err != nil {
		return nil, err
	}

	DiscardWhitespace(r)

	bR, err := r.PeekRune()
	if err != nil {
		return nil, err
	}
	if bR == 'R' {
		_, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}
		return &IndirectRef{i, b}, nil
	}

	rs, err := r.ReadString(3)
	if err != nil {
		return nil, err
	}

	if rs != "obj" {
		return nil, &BadIndirectSpecifierError{r.Position(), "obj", rs}
	}

	DiscardWhitespace(r) // whitespace after obj

	o := IndirectVal{ID: i}
	// FIXME this can also be an indirect object
	o.value, err = ReadDirectObject(r, b)
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

	return &o, nil
}

func (i *IndirectVal) Write(w file.Writer) error {
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

	err = i.Value().Write(w)
	if err != nil {
		return err
	}

	_, err = w.WriteString("\r\nendobj")

	return err
}

func (i *IndirectVal) String() string {
	return fmt.Sprintf("indirect(%s):%s", i.ID.String(), i.Value().String())
}

func (i *IndirectVal) Value() shared.Object {
	return i.value
}
func (i *IndirectRef) Write(w file.Writer) error {
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

	_, err = w.WriteString(" R")

	return err
}

func (i *IndirectRef) String() string {
	return fmt.Sprintf("indirectRef(%s)", i.ID.String())
}

func (i *IndirectRef) Value() shared.Object {
	return i.body.Resolve(&i.ID)
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
