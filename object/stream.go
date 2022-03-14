package object

import (
	"fmt"
	"strconv"

	"github.com/scrouthtv/go-pdf/file"
)

// Stream is a (virtually unlimited) sequence of bytes.
type Stream struct {
	BlobS uint64
	BlobE uint64
	Dict  *Dict
}

func HasStream(r file.Reader) bool {
	DiscardWhitespace(r)

	s, err := r.PeekString(6)
	if err != nil {
		return false
	}

	return s == "stream"
}

func ReadStream(r file.Reader, d *Dict) (*Stream, error) {
	t, err := r.ReadString(6)
	if err != nil {
		return nil, err
	}

	if t != "stream" {
		return nil, &BadStreamStartError{r.Position(), t}
	}

	err = DiscardEOL(r)
	if err != nil {
		return nil, err
	}

	s := Stream{BlobS: uint64(r.Position()), Dict: d}

	//TODO: Since PDF 1.2 content can be in another file

	lenobj, ok := d.Dict["Length"]
	if ok {
		switch v := lenobj.(type) {
		case Integer:
			// dict#length is an integer, use it to read
			return readStreamWithLength(r, &s, int(v))
		case Indirect:
			val := v.Value()
			len, ok := val.(Integer)
			if ok {
				// dict#length is a reference to an integer, use that to read
				return readStreamWithLength(r, &s, int(len))
			} else {
				// dict#length is a reference to something but an integer, read blindly
				return readStreamBlind(r, &s)
			}
		default:
			// dict#length is not an integer, read blindly
			return readStreamBlind(r, &s)
		}
	} else {
		// dict does not have #length, read blindly
		return readStreamBlind(r, &s)
	}
}

func readStreamWithLength(r file.Reader, target *Stream, length int) (*Stream, error) {
	err := r.Advance(length)
	if err != nil {
		return nil, err
	}

	target.BlobE = uint64(r.Position())
	err = DiscardEOL(r)
	if err != nil {
		return nil, err
	}

	e, err := r.ReadString(9)
	if err != nil {
		return nil, err
	}

	if e != "endstream" {
		return nil, &LateEndOfDataMarkerError{r.Position()}
	}

	return target, nil
}

func readStreamBlind(r file.Reader, target *Stream) (*Stream, error) {
	panic("todo")
}

func (s *Stream) String() string {
	return "Stream(" + strconv.FormatUint(s.BlobS, 10) + "-" + strconv.FormatUint(s.BlobE, 10) + ")"
}

func (s *Stream) Write(r file.Writer) error {
	panic("NOT IMPLEMENTED!") //TODO
}

type BadStreamStartError struct {
	Pos   int
	Start string
}

func (e *BadStreamStartError) Error() string {
	return fmt.Sprintf("expected stream at pos %d, got %q instead", e.Pos, e.Start)
}

type EarlyEndOfDataMarkerError struct {
	Pos int
}

func (e *EarlyEndOfDataMarkerError) Error() string {
	return fmt.Sprintf("early stream end of data marker at pos %d", e.Pos)
}

type LateEndOfDataMarkerError struct {
	Pos int
}

func (e *LateEndOfDataMarkerError) Error() string {
	return fmt.Sprintf("late stream end of data marker, expected at pos %d", e.Pos)
}

type ReadingStreamDataError struct {
	Len int
	Err error
}

func (e *ReadingStreamDataError) Error() string {
	return fmt.Sprintf("encountered an error reading the stream's data with %d bytes: %s", e.Len, e.Err.Error())
}

func (e *ReadingStreamDataError) Unwrap() error {
	return e.Err
}
