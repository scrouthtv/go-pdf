package object

import (
	"fmt"

	"github.com/scrouthtv/go-pdf/file"
)

// Dict represents a pdf dictionary.
type Dict struct {
	// Dict is the actual dictionary.
	// Using the builtin map is fine, as
	//  - the keys are unique
	//  - the ordering does not matter
	Dict map[Name]Object
}

func ReadDict(r file.Reader) (*Dict, error) {
	// Dictionary entry with value null shall be ignored.
	// Keys are unique, as names are unique.
	// A dict may contain 0 entries.
	// For using a stream object as value, it should be referenced
	// using an indirect object.
	start, err := r.ReadString(2)
	if err != nil {
		return nil, err // TODO pack error
	}

	if start != "<<" {
		return nil, &BadDictStartError{r.Position(), start}
	}

	end, err := r.PeekString(2)
	if err != nil {
		return nil, err // TODO pack error
	}

	d := Dict{make(map[Name]Object)}

	for end != ">>" {
		k, err := ReadName(r)
		if err != nil {
			return nil, err // TODO pack error
		}

		v, err := ReadObject(r)
		if err != nil {
			return nil, err // TODO pack error
		}

		d.Dict[*k] = v

		end, err = r.PeekString(2)
		if err != nil {
			return nil, err // TODO pack error
		}
	}

	// end is detected by the closing characters,
	// no need to unread anything

	return &d, nil
}

func (d *Dict) Write(w file.Writer) error {
	panic("not impl")
}

func (d *Dict) String() string {
	panic("not impl")
}

type BadDictStartError struct {
	Pos   int
	Start string
}

func (e *BadDictStartError) Error() string {
	return fmt.Sprintf("expected dictionary at pos %d, got %s instead", e.Pos, e.Start)
}
