package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	cg "tango/src/tango/codegen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file name>\n", os.Args[0])
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read file: %s\nError: %v\n", os.Args[1], err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var code []cg.IRIns

	for scanner.Scan() {
		line := scanner.Text()

		// Get Label
		label := cg.IRLabel("")
		colonIndex := strings.Index(line, ":")
		if colonIndex != -1 {
			label = cg.IRLabel(line[:colonIndex])
		}
		line = line[colonIndex+1:]

		// Split into various parts
		splitted := strings.Split(line, ",")

		// Get op code
		op := cg.IROp(splitted[0])

		// Get type
		typ := cg.GetType(op)
		if typ == cg.INV {
			log.Fatal("Invalid operator: " + op)
		}

		// Get args
		arg1, arg2, dst := cg.GetRegs(splitted, typ, op)

		ins := cg.IRIns{
			Typ:   typ,
			Op:    op,
			Arg1:  arg1,
			Arg2:  arg2,
			Dst:   dst,
			Label: label,
		}

		code = append(code, ins)
	}

	cg.GenBBLList(code)

	fmt.Printf("%v\n", cg.BBLList)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
