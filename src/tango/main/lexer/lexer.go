package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"tango/lexer"
	"tango/token"
)

func main() {
	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Unable to read test file: 1.go\n")
		return
	}
	l := lexer.NewLexer([]byte(input))
	for tok := l.Scan(); tok.Type != token.EOF; tok = l.Scan() {
		switch {
		case tok.Type == token.INVALID:
			fmt.Printf("Invalid token found: %s\n", tok.Lit)
			return
		default:
			fmt.Printf("Token found: %s\n", tok.Lit)
		}
	}
}
