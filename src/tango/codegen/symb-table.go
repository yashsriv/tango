package codegen

import (
	"log"
	"strconv"
)

// SymbolTableEntry is an entry in the SymbolTable
type SymbolTableEntry interface {
}

// SymbolTableLiteralEntry refers to a literal in the symbol table
type SymbolTableLiteralEntry struct {
	Repr  string
	Value int
}

// SymbolTableTargetEntry refers to a target in the symbol table
type SymbolTableTargetEntry struct {
	Target string
}

// SymbolTableRegisterEntry refers to a register in the symbol table
type SymbolTableRegisterEntry struct {
	Register string
}

// SymbolTable is an array of SymbolTableEntries
var SymbolTable []SymbolTableEntry
var symbolMap = make(map[string]SymbolTableEntry)

func insertToSymbolTable(val string) SymbolTableEntry {
	var entry SymbolTableEntry
	switch val[0] {
	case '$':
		i, err := strconv.ParseInt(val[1:], 0, 32)
		if err != nil {
			log.Fatalf("Unable to parse int: %d\nError: %v\n", i, err)
		}

		entry = &SymbolTableLiteralEntry{
			Repr:  val,
			Value: int(i),
		}
	case '#':
		entry = &SymbolTableTargetEntry{
			Target: val[1:],
		}
	case 'r':
		entry = &SymbolTableRegisterEntry{
			Register: val,
		}
	}
	SymbolTable = append(SymbolTable, entry)
	return entry
}

// GetRegs populates symbol table and gets virtual registers
func GetRegs(splitted []string, typ IRType, op IROp) (SymbolTableEntry, SymbolTableEntry, SymbolTableEntry) {

	var arg1, arg2, dst SymbolTableEntry
	if typ == BOP || typ == CBR {
		if val, ok := symbolMap[splitted[1]]; ok {
			dst = val
		} else {
			dst = insertToSymbolTable(splitted[1])
			symbolMap[splitted[1]] = dst
		}
		if val, ok := symbolMap[splitted[2]]; ok {
			arg1 = val
		} else {
			arg1 = insertToSymbolTable(splitted[2])
			symbolMap[splitted[1]] = arg1
		}
		if val, ok := symbolMap[splitted[3]]; ok {
			arg2 = val
		} else {
			arg2 = insertToSymbolTable(splitted[3])
			symbolMap[splitted[1]] = arg2
		}
	} else if typ == UOP || typ == ASN {
		if val, ok := symbolMap[splitted[1]]; ok {
			dst = val
		} else {
			dst = insertToSymbolTable(splitted[1])
			symbolMap[splitted[1]] = dst
		}
		if val, ok := symbolMap[splitted[2]]; ok {
			arg1 = val
		} else {
			arg1 = insertToSymbolTable(splitted[2])
			symbolMap[splitted[1]] = arg1
		}
	} else if typ == JMP {
		if val, ok := symbolMap[splitted[1]]; ok {
			arg1 = val
		} else {
			arg1 = insertToSymbolTable(splitted[1])
			symbolMap[splitted[1]] = arg1
		}
	} else if typ == KEY {
		if !(op == RET || op == HALT) {
			if val, ok := symbolMap[splitted[1]]; ok {
				arg1 = val
			} else {
				arg1 = insertToSymbolTable(splitted[1])
				symbolMap[splitted[1]] = arg1
			}
		}
	}
	return arg1, arg2, dst
}
