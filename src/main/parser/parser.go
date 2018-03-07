package main

import (
	"fmt"
	"log"
	"os"

	"tango/src/ast"
	"tango/src/lexer"
	"tango/src/parser"
	"tango/src/token"
)

func reverseSlice(input ast.Stack) ast.Stack {
	if len(input) == 0 {
		return input
	}
	return append(reverseSlice(input[1:]), input[0])
}

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
	if !ok {
		panic("Expected a Source File")
	}

	outFile, err := os.Create("file1.html")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// sourceFile.GenOutput()

	// fmt.Println(ast.Nodes)

	prefix := make(ast.Stack, 0)
	suffix := make(ast.Stack, 0)

	prev := sourceFile
	derivation := ast.Derivations[sourceFile]
	for {
		// fmt.Println(prev, prefix, derivation, reverseSlice(suffix))
		revsuffix := reverseSlice(suffix)
		fmt.Printf("%s _%s_ %s => %s %s %s\n", prefix, prev, revsuffix, prefix, derivation, revsuffix)
		found := false
		var next []ast.Attrib
		for i := len(derivation) - 1; i >= 0; i-- {
			switch v := derivation[i].(type) {
			case *token.Token:
				if !found {
					suffix = suffix.Push(v)
				} else {
					prefix = prefix.Push(v)
				}
			case *ast.Node:
				if !found {
					found = true
					prev = v
					next = ast.Derivations[v]
				} else {
					prefix = prefix.Push(v)
				}
			default:
				log.Fatalf("%T\n", v)
			}
		}
		if !found {
			for !prefix.Empty() {
				var attrib ast.Attrib
				prefix, attrib = prefix.Pop()

				switch v := attrib.(type) {
				case *token.Token:
					suffix = suffix.Push(v)
				case *ast.Node:
					found = true
					prev = v
					next = ast.Derivations[v]
					goto done
				}
			}
		}
	done:
		if !found {
			break
		}
		derivation = next
	}
	fmt.Println(reverseSlice(suffix))
}