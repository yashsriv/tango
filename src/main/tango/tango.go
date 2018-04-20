package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"tango/src/ast"
	"tango/src/codegen"
	"tango/src/lexer"
	"tango/src/parser"
)

func main() {

	langPtr := flag.String("x", "None", "Language to output. Possible values are IR, ASM and None")
	outFile := flag.String("o", "./a.out", "Filename of output executable. Only valid with language None")

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s [<options>] <file name>\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}
	lex, err := lexer.NewWrapperFile(flag.Args()[0])
	if err != nil {
		fmt.Printf("Unable to read file: %s\n", flag.Args()[0])
		return
	}
	p := parser.NewParser()
	st, err := p.Parse(lex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: An error occurred while parsing: %v\n", err)
		os.Exit(1)
	}
	sourceFile, ok := st.(*ast.AddrCode)
	if !ok || !sourceFile.TopLevel {
		fmt.Fprintln(os.Stderr, "Error: Unexpected end of file")
		os.Exit(1)
	}

	if *langPtr != "None" && *outFile != "./a.out" {
		fmt.Fprintln(os.Stderr, "Warning: -o flag is only valid when -x is None (the default)")
	}

	switch *langPtr {
	case "IR":
		for _, val := range sourceFile.Code {
			fmt.Println(val)
		}
	case "ASM", "None":
		codegen.GenBBLList(sourceFile.Code)
		codegen.GenerateASM()
		if *langPtr == "ASM" {
			fmt.Println(codegen.Code)
			break
		}
		gccCmd := exec.Command("gcc", "-m32", "-o", *outFile, "-x", "assembler", "-")
		gccIn, err := gccCmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		gccOut, err := gccCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		gccErr, err := gccCmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		gccCmd.Start()
		gccIn.Write([]byte(codegen.Code))
		gccIn.Close()

		_, err = io.Copy(os.Stdout, gccOut)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(os.Stderr, gccErr)
		if err != nil {
			panic(err)
		}
		err = gccCmd.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "gcc closed with an error. Assembly code written to %s\n", *outFile+".s")
			outFile, _ := os.Create(*outFile + ".s")
			fmt.Fprintln(outFile, codegen.Code)
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported -x param: %s\n", *langPtr)
		os.Exit(1)

	}

}
