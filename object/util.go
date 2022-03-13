package object

import "github.com/scrouthtv/go-pdf/file"

func readHexCharacter(r file.Reader) (rune, error) {
	panic("not impl")
	// TODO
}

func isTokenDelimiter(r rune) bool {
	// Are there any other token delimiters??
	return r == ' ' || r == '(' || r == ')' ||
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
