package object_test

var test_strings = map[string]string{
	"(This is a string)":                        "This is a string",
	"(contain\r\nnewlines and such.)":           "contain\r\nnewlines and such.",
	"(contain\nnewlines and such.)":             "contain\nnewlines and such.",
	"(balanced parens ( ) )":                    "balanced parens ( )",
	"(special chars ( * ! & } ^ %and so on) .)": "special chars ( * ! & } ^ %and so on) .",
	"()": "",
}

/*(These \
two strings \
are the same.)
(These two strings are the same.)

(This string has an end-of-line at the end of it.
)
(So does this one.\n)

(This string contains \245two octal characters\307.)

the literal
(\0053)
denotes a string containing two characters, \005 (Control-E) followed by the digit 3, whereas both
(\053)
and
(\53)
denote strings containing the single character \053, a plus sign (+).
}*/
