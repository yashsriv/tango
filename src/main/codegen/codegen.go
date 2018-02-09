package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	cg "tango/src/codegen"
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

		if strings.TrimSpace(line) == "" {
			continue
		}

		// Get Label
		colonIndex := strings.Index(line, ":")
		if colonIndex != -1 {
			label := line[:colonIndex]
			var dst cg.SymbolTableEntry
			dst = cg.InsertToSymbolTable(label)
			ins := cg.IRIns{
				Typ: cg.LBL,
				Dst: dst,
			}
			code = append(code, ins)
			continue
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
		arg1, arg2, dst, err := cg.GetRegs(splitted, typ, op)
		if err != nil {
			log.Fatalf("Error file parsing args for: %s\nError: %v\n", line, err)
		}

		ins := cg.IRIns{
			Typ:  typ,
			Op:   op,
			Arg1: arg1,
			Arg2: arg2,
			Dst:  dst,
		}

		code = append(code, ins)
	}

	cg.GenBBLList(code)

	for _, b := range cg.BBLList {
		fmt.Print(b.String())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
