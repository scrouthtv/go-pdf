package object

import (
	"go-pdf/pdfio"
	"strconv"
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
