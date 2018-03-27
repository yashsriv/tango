package main

import (
	"fmt"
	"os"

	"tango/src/ast"
	"tango/src/lexer"
	"tango/src/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file name>\n", os.Args[0])
		os.Exit(1)
	}
	lex, err := lexer.NewWrapperFile(os.Args[1])
	if err != nil {
		fmt.Printf("Unable to read file: %s\n", os.Args[1])
		return
	}
	p := parser.NewParser()
	st, err := p.Parse(lex)
	if err != nil {
		panic(err)
	}
	sourceFile, ok := st.(*ast.AddrCode)
	if !ok {
		panic("Expected a AddrCode")
	}

	for _, val := range sourceFile.Code {
		fmt.Println(val)
	}
}
