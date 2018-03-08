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

type entry struct {
	prefix     ast.Stack
	suffix     ast.Stack
	node       *ast.Node
	derivation ast.Stack
}

func genOutput(sourceFile *ast.Node) {
	prefix := make(ast.Stack, 0)
	suffix := make(ast.Stack, 0)

	entries := make([]entry, 0)

	prev := sourceFile
	found := true
	for found {
		revsuffix := reverseSlice(suffix)
		entries = append(entries, entry{
			prefix:     prefix,
			node:       prev,
			suffix:     revsuffix,
			derivation: ast.Derivations[prev],
		})
		found, prev, prefix, suffix = findNext(prefix, prev, suffix)
	}
	entries = append(entries, entry{
		prefix:     prefix,
		node:       nil,
		suffix:     reverseSlice(suffix),
		derivation: ast.Stack{},
	})

	outputFormatting(entries)
}

func outputFormatting(entries []entry) {
	// TODO: Load a template html file and populate it
	// See https://astaxie.gitbooks.io/build-web-application-with-golang/en/07.4.html
	// _, filename, _, _ := runtime.Caller(1)
	// f, err := os.Open(path.Join(path.Dir(filename), "data.csv"))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(filepath.Join(filepath.Dir(filename), "templ.html"))

	for _, val := range entries {
		fmt.Printf("%s _%s_ %s\n", val.prefix, val.node, val.suffix)
	}
	fmt.Println("$")
}

func findNext(prefix ast.Stack, prev *ast.Node, suffix ast.Stack) (bool, *ast.Node, ast.Stack, ast.Stack) {
	derivation := ast.Derivations[prev]
	found := false
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
			} else {
				prefix = prefix.Push(v)
			}
		default:
			log.Fatalf("Unknown type: %T\n", v)
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
				goto done
			default:
				log.Fatalf("Unknown type: %T\n", v)
			}
		}
	}
done:
	return found, prev, prefix, suffix
}

func reverseSlice(input ast.Stack) ast.Stack {
	if len(input) == 0 {
		return input
	}
	return append(reverseSlice(input[1:]), input[0])
}
