package codegen

import (
	"errors"
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

var symbolMap = make(map[string]SymbolTableEntry)

// InsertToSymbolTable inserts a single entry into table
func InsertToSymbolTable(val string) SymbolTableEntry {
	if val, ok := symbolMap[val]; ok {
		return val
	}
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
	symbolMap[val] = entry
	return entry
}

// GetRegs populates symbol table and gets virtual registers
func GetRegs(splitted []string, typ IRType, op IROp) (SymbolTableEntry, SymbolTableEntry, SymbolTableEntry, error) {

	var arg1, arg2, dst SymbolTableEntry

	if typ == LBL {
		return nil, nil, nil, errors.New("We Should never GetRegs for a label")
	}

	if typ == BOP || typ == CBR {
		if len(splitted) < 4 {
			return nil, nil, nil, errors.New("Not enough args to a binary operand")
		}
		dst = InsertToSymbolTable(splitted[1])
		arg1 = InsertToSymbolTable(splitted[2])
		arg2 = InsertToSymbolTable(splitted[3])
	} else if typ == UOP || typ == ASN {
		if len(splitted) < 3 {
			return nil, nil, nil, errors.New("Not enough args to a unary operand")
		}
		dst = InsertToSymbolTable(splitted[1])
		arg1 = InsertToSymbolTable(splitted[2])
	} else if typ == JMP {
		if len(splitted) < 2 {
			return nil, nil, nil, errors.New("Not enough args to a jump operand")
		}
		arg1 = InsertToSymbolTable(splitted[1])
	} else if typ == KEY {
		if !(op == RET || op == HALT) {
			if len(splitted) < 2 {
				return nil, nil, nil, errors.New("Not enough args to a call/key operand")
			}
			arg1 = InsertToSymbolTable(splitted[1])
		}
	}
	return arg1, arg2, dst, nil
}
