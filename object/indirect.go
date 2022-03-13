package object

import (
	"fmt"
	"io"

	"github.com/scrouthtv/go-pdf/file"
)

type ObjId struct {
	// Id is the object number. Object numbers can be arbitrarily assigned.
	// However, they should be unique.
	Id Integer

	// Gen is the object's generation number. It is non-negative and
	// defaults to 0. Other generation numbers only appear in later file
	// updates.
	Gen Integer
}

func (o *ObjId) String() string {
	return fmt.Sprintf("%d/%d", o.Id, o.Gen)
}

type Indirect struct {
	Id    ObjId
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
	s, err := r.ReadString(5)
	if err != nil {
		return false
	}

	return s == "obj\r\n"
}

func ReadIndirect(r file.Reader) (*Indirect, error) {
	i := Indirect{}
	var err error

	i.Id.Id, err = ReadInteger(r)
	if err != nil {
		return nil, err
	}

	read, _, err := r.ReadRune()
	if err != nil {
		return nil, err
	}

	if read != ' ' {
		return nil, &ErrRunawayIndirectMember{r.Position()}
	}

	i.Id.Gen, err = ReadInteger(r)
	if err != nil {
		return nil, err
	}

	if read != ' ' {
		return nil, &ErrRunawayIndirectMember{r.Position()}
	}

	// FIXME this can also be an indirect object
	i.Value, err = ReadDirectObject(r)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (i *Indirect) Write(w file.Writer) error {
	err := i.Id.Id.Write(w)
	if err != nil {
		return err
	}

	err = w.WriteRune(' ')
	if err != nil {
		return err
	}

	err = i.Id.Gen.Write(w)
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
	panic("not impl")
}

type ErrRunawayIndirectMember struct {
	Position int
}

func (e *ErrRunawayIndirectMember) Error() string {
	return fmt.Sprintf("runaway indirect member, excepted space at %d", e.Position)
}