package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/awalterschulze/gographviz"

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
	sourceFile, ok := st.(*ast.SourceFile)
	if !ok {
		panic("Expected a Source File")
	}

	outFile, err := os.Create("file1.html")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Generate dot graph
	graphAst, _ := gographviz.ParseString(`digraph main {}`)
	graph := gographviz.NewGraph()
	err = gographviz.Analyse(graphAst, graph)
	if err != nil {
		panic(err)
	}
	sourceFile.GenGraph(graph)

	// Get svg output
	dotCmd := exec.Command("dot", "-Tsvg")

	// Input pipe to running dot process
	dotIn, err := dotCmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	// Output pipe from running dot process
	dotOut, err := dotCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// Start running dot process
	err = dotCmd.Start()
	if err != nil {
		panic(err)
	}

	// Write to dot process
	dotIn.Write([]byte(graph.String()))
	dotIn.Close()

	_, err = io.Copy(outFile, dotOut)
	if err != nil {
		panic(err)
	}

	defer dotCmd.Wait()

}
