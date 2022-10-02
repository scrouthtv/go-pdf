package object

import (
	"fmt"
	"io"
	"strconv"

	"go-pdf/pdfio"
)

// Stream is a (virtually unlimited) sequence of bytes.
type Stream struct {
	BlobS uint64
	BlobE uint64
	Dict  *Dict
}

func HasStream(r pdfio.Reader) bool {
	pdfio.DiscardWhitespace(r)

	s, err := r.PeekString(6)
	if err != nil {
		return false
	}

	return s == "stream"
}

func ReadStream(r pdfio.Reader, d *Dict) (*Stream, error) {
	t, err := r.ReadString(6)
	if err != nil {
		return nil, err
	}

	if t != "stream" {
		return nil, &BadStreamStartError{r.Position(), t}
	}

	err = pdfio.DiscardEOL(r)
	if err != nil {
		return nil, err
	}

	s := Stream{BlobS: uint64(r.Position()), Dict: d}

	//TODO: Since PDF 1.2 content can be in another file

	println("ok")

	lenobj, ok := d.Dict["Length"]
	if ok {
		switch v := lenobj.(type) {
		case Integer:
			// dict#length is an integer, use it to read
			return readStreamWithLength(r, &s, int(v))
		case Indirect:
			val := v.Value()
			len, ok := val.(Integer)
			println(v.String())
			if ok {
				println("b")
				// dict#length is a reference to an integer, use that to read
				return readStreamWithLength(r, &s, int(len))
			} else {
				println("c")
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

// readStreamWithLength takes a stream object
// that is already populated with start position and dict
// and sets the end position to use the specified length.
//
// End tokens inside the stream are not checked.
// If the stream does not end with "endstream" (EOL markers
// are ignored), the method returns an error.
func readStreamWithLength(r pdfio.Reader, target *Stream, length int) (*Stream, error) {
	println("read w length")
	_, err := r.Seek(int64(length), io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	target.BlobE = uint64(r.Position())
	err = pdfio.DiscardEOL(r)
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

// readStreamBlind takes a stream object that is already
// populated with start position and dict and sets the
// end position.
//
// The end position is determined by searching for an EOL marker,
// that is followed by "endstream".
// The end position is set to the position of that EOL marker.
func readStreamBlind(r pdfio.Reader, target *Stream) (*Stream, error) {
	for {
		read, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}

		if pdfio.IsEOL(read) {
			target.BlobE = uint64(r.Position() - 1)
			pdfio.DiscardEOL(r) // discard additional LF

			reads, err := r.ReadString(9)
			if err != nil {
				return nil, err
			}

			if reads == "endstream" {
				return target, nil
			}
		}
	}
}

func (s *Stream) String() string {
	return "Stream(" + strconv.FormatUint(s.BlobS, 10) + "-" + strconv.FormatUint(s.BlobE, 10) + ")"
}

func (s *Stream) Write(r pdfio.Writer) error {
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
