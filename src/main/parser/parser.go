package main

import (
	"fmt"
	"os"

	"tango/src/ast"
	"tango/src/html"
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
	sourceFile, ok := st.(*ast.Node)
	if !ok || sourceFile.String() != "SourceFile" {
		panic("Expected a Source File")
	}

	genOutput(sourceFile)
}

func genOutput(sourceFile *ast.Node) {

	entries := make([]html.Entry, 0)

	derivations := ast.Stack{sourceFile}
	found := true
	// prev := sourceFile
	for found {
		// fmt.Println("["+prefix.String()+"]", prev, "["+suffix.String()+"]")
		var index int
		found, index = findNext(derivations)
		if found {
			entries = append(entries, html.Entry{
				Prefix: derivations[:index],
				Node:   derivations[index].(*ast.Node),
				Suffix: derivations[index+1:],
			})
			newderivations := make(ast.Stack, 0)
			for _, v := range derivations[:index] {
				newderivations = append(newderivations, v)
			}
			for _, v := range ast.Derivations[derivations[index].(*ast.Node)] {
				newderivations = append(newderivations, v)
			}
			for _, v := range derivations[index+1:] {
				newderivations = append(newderivations, v)
			}
			derivations = newderivations
		} else {
			entries = append(entries, html.Entry{
				Prefix: derivations,
				Node:   nil,
				Suffix: ast.Stack{},
			})
		}
	}

	f, err := os.Create("./file1.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	html.Output(entries, f)

	f.Sync()

}

func findNext(derivations ast.Stack) (bool, int) {
	for i := len(derivations) - 1; i >= 0; i-- {
		switch derivations[i].(type) {
		case *ast.Node:
			return true, i
		}
	}
	return false, -1
}

func reverseSlice(input ast.Stack) ast.Stack {
	if len(input) == 0 {
		return input
	}
	return append(reverseSlice(input[1:]), input[0])
}
