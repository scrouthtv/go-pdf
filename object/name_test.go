package object_test

var test_names = map[string]string{
	"/Name1":                             "Name1",
	"/ASomewhatLongerName":               "ASomewhatLongerName",
	"/A;Name_With-Various***Characters?": "A;Name_With-Various***Characters?",
	"/1.2":                               "1.2",
	"/$$":                                "$$",
	"/@pattern":                          "@pattern",
	"/.notdef":                           ".notdef",
	"/Lime#20Green":                      "Lime Green",
	"/paired#28#29parentheses":           "paired()parentheses",
	"/The_Key_of_F#23_Minor":             "The_Key_of_F#_Minor",
	"/A#42":                              "AB",
}
