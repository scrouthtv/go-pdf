package object

import (
	"fmt"
	"strconv"

	"github.com/scrouthtv/go-pdf/file"
)

type String struct {
	S     string
	isHex bool
}

func ReadString(r file.Reader) (*String, error) {
	id, _, err := r.ReadRune()
	if err != nil {
		return nil, err // TODO pack error
	}

	switch id {
	case '(':
		return readLiteralString(r)
	case '<':
		return readHexString(r)
	default:
		return nil, &BadStringStartError{r.Position(), id}
	}
}

func readLiteralString(r file.Reader) (*String, error) {
	parens := 0
	out := String{"", false}

	for {
		read, _, err := r.ReadRune()
		if err != nil {
			return nil, &RunawayStringError{r.Position(), ')', err}
		}

		if read == '\\' {
			read, _, err = r.ReadRune()
			if err != nil {
				return nil, &RunawayEscapeError{r.Position(), err}
			}

			switch read {
			case 'n':
				out.S += "\n"
			case 'r':
				out.S += "\r"
			case 't':
				out.S += "\t"
			case 'b':
				out.S += "\b"
			case 'f':
				out.S += "\f"
			case '(':
				out.S += "("
			case ')':
				out.S += ")"
			case '\\':
				out.S += "\\"
			case '\r':
				// discard backslash and EOL marker when escaped.
				// EOL marker always consists of CRLF.
				read, _, err = r.ReadRune()
				if err != nil {
					return nil, err // TODO pack error
				}

				if read != '\n' {
					return nil, &RunawayEscapeError{r.Position(), err}
				}
			default:
				read, err = readOctal(r, read)
				if err != nil {
					return nil, err
				}

				out.S += string(read)
			}

			continue
		} else if read == '(' {
			parens++
		} else if read == ')' {
			parens--
			if parens < 0 {
				// end is detected by the closing character,
				// so we don't need to unread it.
				return &out, nil
			}
		}

		out.S += string(read)
	}
}

func readOctal(r file.Reader, pre rune) (rune, error) { // TODO The number ddd may consist of one, two, or three octal digits
	s, err := r.ReadString(2)
	if err != nil {
		return -1, err // TODO pack error
	}

	res, err := strconv.ParseUint(string(pre)+s, 8, 32)

	return rune(res), err
}

// readHexString reads a string of consecutive two-digit integers
// and parses them as hexadecimal numbers.
// It ends when reaching a '>'.
func readHexString(r file.Reader) (*String, error) {
	out := String{"", false}

	for {
		next, err := r.PeekRune()
		if err != nil {
			return nil, err
		}
		if next == '>' {
			r.ReadRune()
			break
		}
		c, err := readHexCharacter(r)
		if err != nil {
			return nil, err
		}
		out.S += string(c)
	}
	return &out, nil
}

func (s *String) Write(w file.Writer) error {
	if s.isHex {
		panic("not impl")
	} else {
		_, err := w.Write([]byte{'('})
		if err != nil {
			return err // TODO pack error
		}

		_, err = w.WriteString(s.S)
		if err != nil {
			return err // TODO pack error
		}

		_, err = w.Write([]byte{')'})
		return err
	}
}

func (s *String) String() string {
	return "string(" + s.S + ")"
}

type BadStringStartError struct {
	Pos   int
	Start rune
}

func (e *BadStringStartError) Error() string {
	return fmt.Sprintf("expected string at pos %d, got %q instead", e.Pos, e.Start)
}

type RunawayStringError struct {
	Pos    int
	Closer rune
	Err    error
}

func (e *RunawayStringError) Error() string {
	return fmt.Sprintf("got error %s reading string at %d while looking for %q", e.Err, e.Pos, e.Closer)
}

func (e *RunawayStringError) Unwrap() error {
	return e.Err
}

type RunawayEscapeError struct {
	Pos int
	Err error
}

func (e *RunawayEscapeError) Error() string {
	return fmt.Sprintf("got error %s reading escape sequence at %d", e.Err, e.Pos)
}

func (e *RunawayEscapeError) Unwrap() error {
	return e.Err
}
