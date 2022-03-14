package object

import (
	"fmt"

	"github.com/scrouthtv/go-pdf/file"
)

// Stream is a (virtually unlimited) sequence of bytes.
type Stream struct {
	BlobS uint64
	BlobE uint64
	Dict  *Dict
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

	s := &Stream{BlobS: uint64(r.Position()), Dict: d}

	//TODO: Since PDF 1.2 content can be in extra file

	err = r.Advance(int(d.Dict["Length"].(Integer)))
	if err != nil {
		return nil, err
	}
	s.BlobE = uint64(r.Position())
	err = DiscardEOL(r)
	if err != nil {
		return nil, err
	}

	e, err := r.ReadString(9)

	if e != "endstream" {
		return nil, &LateEndOfDataMarkerError{r.Position()}
	}

	return s, nil
}

func (s *Stream) String() string {
	panic("NOT IMPLEMENTED!") //TODO
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
