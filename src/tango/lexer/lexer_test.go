package lexer

import (
	"io/ioutil"
	"testing"

	"tango/token"
)

var files = []string{
	"../../../test/1.go",
	"../../../test/2.go",
	"../../../test/3.go",
	"../../../test/4.go",
}

func Test(t *testing.T) {
	for _, x := range files {
		input, err := ioutil.ReadFile(x)
		if err != nil {
			t.Errorf("Unable to read test file: %s", x)
		}
		l := NewWrapper(input)
		for tok := l.Scan(); tok.Type != token.EOF; tok = l.Scan() {
			switch {
			case tok.Type == token.INVALID:
				t.Errorf("Invalid token found: %s", tok.Lit)
			}
		}
	}
}
