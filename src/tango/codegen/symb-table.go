package codegen

import (
	"fmt"
	"log"
	"strconv"
)

// SymbolTableEntry is an entry in the SymbolTable
type SymbolTableEntry interface {
	SymbolTableString() string
}

// SymbolTableLiteralEntry refers to a literal in the symbol table
type SymbolTableLiteralEntry struct {
	Repr  string
	Value int
}

// SymbolTableString returns a string representation and also ensures types
func (s SymbolTableLiteralEntry) SymbolTableString() string {
	return fmt.Sprintf("$%d", s.Value)
}

// SymbolTableTargetEntry refers to a target in the symbol table
type SymbolTableTargetEntry struct {
	Target string
}

// SymbolTableString returns a string representation and also ensures types
func (s SymbolTableTargetEntry) SymbolTableString() string {
	return fmt.Sprintf("#%s", s.Target)
}

// SymbolTableRegisterEntry refers to a register in the symbol table
type SymbolTableRegisterEntry struct {
	Register string
}

// SymbolTableString returns a string representation and also ensures types
func (s SymbolTableRegisterEntry) SymbolTableString() string {
	return fmt.Sprintf("%%%s", s.Register)
}

// SymbolTable is an array of SymbolTableEntries
var SymbolTable []SymbolTableEntry
var SymbolMap = make(map[string]SymbolTableEntry)

func InsertToSymbolTable(val string) SymbolTableEntry {
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
	default:
		entry = &SymbolTableTargetEntry{
			Target: val,
		}
	}
	SymbolTable = append(SymbolTable, entry)
	return entry
}

// GetRegs populates symbol table and gets virtual registers
func GetRegs(splitted []string, typ IRType, op IROp) (SymbolTableEntry, SymbolTableEntry, SymbolTableEntry) {

	var arg1, arg2, dst SymbolTableEntry
	if typ == BOP || typ == CBR {
		if val, ok := SymbolMap[splitted[1]]; ok {
			dst = val
		} else {
			dst = InsertToSymbolTable(splitted[1])
			SymbolMap[splitted[1]] = dst
		}
		if val, ok := SymbolMap[splitted[2]]; ok {
			arg1 = val
		} else {
			arg1 = InsertToSymbolTable(splitted[2])
			SymbolMap[splitted[2]] = arg1
		}
		if val, ok := SymbolMap[splitted[3]]; ok {
			arg2 = val
		} else {
			arg2 = InsertToSymbolTable(splitted[3])
			SymbolMap[splitted[3]] = arg2
		}
	} else if typ == UOP || typ == ASN {
		if val, ok := SymbolMap[splitted[1]]; ok {
			dst = val
		} else {
			dst = InsertToSymbolTable(splitted[1])
			SymbolMap[splitted[1]] = dst
		}
		if val, ok := SymbolMap[splitted[2]]; ok {
			arg1 = val
		} else {
			arg1 = InsertToSymbolTable(splitted[2])
			SymbolMap[splitted[2]] = arg1
		}
	} else if typ == JMP {
		if val, ok := SymbolMap[splitted[1]]; ok {
			arg1 = val
		} else {
			arg1 = InsertToSymbolTable(splitted[1])
			SymbolMap[splitted[1]] = arg1
		}
	} else if typ == KEY {
		if !(op == RET || op == HALT) {
			if val, ok := SymbolMap[splitted[1]]; ok {
				arg1 = val
			} else {
				arg1 = InsertToSymbolTable(splitted[1])
				SymbolMap[splitted[1]] = arg1
			}
		}
	}
	return arg1, arg2, dst
}
