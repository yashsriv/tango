package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ryanuber/columnize"

	"tango/lexer"
	"tango/token"
)

type countMap map[token.Type]int
type setMap map[token.Type]map[string]bool

func main() {
	counts := make(countMap)
	sets := make(setMap)
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file name>\n", os.Args[0])
		os.Exit(1)
	}
	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Unable to read file: %s\n", os.Args[1])
		return
	}
	l := lexer.NewLexer(input)
	for tok := l.Scan(); tok.Type != token.EOF; tok = l.Scan() {
		switch {
		case tok.Type == token.INVALID:
			fmt.Printf("Invalid token found: %s\n", tok.Lit)
			os.Exit(2)
		default:
			counts[tok.Type]++
			set, ok := sets[tok.Type]
			if !ok {
				set = make(map[string]bool)
			}
			set[string(tok.Lit)] = true
			sets[tok.Type] = set
		}
	}
	lines := make([]string, 1)
	lines[0] = fmt.Sprintf("Token λ Occurrances λ Lexemes")
	for t, count := range counts {
		line1 := fmt.Sprintf("%s λ %d λ ", token.TokMap.Id(t), count)
		line1Done := false
		for lit := range sets[t] {
			if !line1Done {
				line1 += lit
				lines = append(lines, line1)
				line1Done = true
			} else {
				lines = append(lines, fmt.Sprintf("λ λ %s", lit))
			}
		}
	}
	config := columnize.DefaultConfig()
	config.Delim = "λ"
	fmt.Println(columnize.Format(lines, config))
}
