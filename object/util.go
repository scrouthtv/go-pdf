package object

import (
	"strconv"

	"go-pdf/pdfio"
)

func readHexCharacter(r pdfio.Reader) (rune, error) {
	l, err := r.ReadString(2)
	if err != nil {
		return rune(0), err
	}
	c, err := strconv.ParseInt(l, 16, 8)
	if err != nil {
		return rune(0), err
	}
	return rune(c), err
}

func DiscardEOL(r pdfio.Reader) error {

	p, err := r.PeekRune()
	if err != nil {
		return err
	}
	if p == '\r' {
		_, _, err = r.ReadRune()
		if err != nil {
			return err
		}
		p, err = r.PeekRune()
		if err != nil {
			return err
		}
	}
	if p == '\n' {
		_, _, err = r.ReadRune()
		if err != nil {
			return err
		}
	}
	return nil
}

func isEOL(r rune) bool {
	return r == '\r' || r == '\n'
}

func isTokenDelimiter(r rune) bool {
	// Are there any other token delimiters??
	return isWhitespace(r) || r == '(' || r == ')' ||
		r == '<' || r == '>' || r == '[' || r == ']' ||
		r == '/' || r == '%'
}

// isRegularCharacter determines whether the specified rune
// shall be escaped inside a name using #.
func isRegularCharacter(r rune) bool {
	if r == '#' { // A number sign shall be written using its hex code.
		return false
	}

	// Regular characters are inside the range excl mark ! through tilde ~.
	return r >= '!' && r <= '~'
}

// isNumericCharacter determines whether the specified character may be
// part of a numeric (either integer or real).
func isNumericCharacter(r rune) bool {
	return (r >= '0' && r <= '9') || r == '+' || r == '-' || r == '.'
}

func isWhitespace(r rune) bool {
	return r == 0 || r == ' ' || r == '\t' ||
		r == '\r' || r == '\n' || r == 12
}

func DiscardWhitespace(r pdfio.Reader) error {
	read, err := r.PeekRune()
	if err != nil {
		return err
	}

	for isWhitespace(read) {
		_, _, err = r.ReadRune()
		if err != nil {
			return err
		}

		read, err = r.PeekRune()
		if err != nil {
			return err
		}
	}

	return nil
}
