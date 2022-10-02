package pdfio

// DiscardEOL discards a single EOL marker (CR, LF, CRLF)
func DiscardEOL(r Reader) error {
	p, err := r.PeekRune()
	if err != nil {
		return err
	}
	if p == '\r' {
		_, _, err = r.ReadRune()
		if err != nil {
			return err
		}

		// Check if CRLF
		p, err = r.PeekRune()
		if err != nil {
			return err
		}
		if p == '\r' {
			_, _, err = r.ReadRune()
			if err != nil {
				return err
			}
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

func IsEOL(r rune) bool {
	return r == '\r' || r == '\n'
}

func IsTokenDelimiter(r rune) bool {
	// Are there any other token delimiters??
	return isWhitespace(r) || r == '(' || r == ')' ||
		r == '<' || r == '>' || r == '[' || r == ']' ||
		r == '/' || r == '%'
}

func isWhitespace(r rune) bool {
	return r == 0 || r == ' ' || r == '\t' ||
		r == '\r' || r == '\n' || r == 12
}

func DiscardWhitespace(r Reader) error {
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
