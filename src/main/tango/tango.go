package main

import (
	"fmt"
	"os"
	"os/exec"

	"tango/src/ast"
	"tango/src/codegen"
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
		panic("Expected an AddrCode")
	}

	for _, val := range sourceFile.Code {
		fmt.Fprintln(os.Stderr, val)
	}
	codegen.GenBBLList(sourceFile.Code)
	codegen.GenerateASM()

	// fmt.Println(codegen.Code)

	gccCmd := exec.Command("gcc", "-m32", "-o", "prog", "-x", "assembler", "-")

	gccIn, err := gccCmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	gccCmd.Start()
	gccIn.Write([]byte(codegen.Code))
	gccIn.Close()
	gccCmd.Wait()

}
