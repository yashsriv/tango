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

// SymbolTableVariableEntry refers to a register in the symbol table
type SymbolTableVariableEntry struct {
	MemoryLocation string
	Declared       bool
	Assignments    int
}

// SymbolTableString returns a string representation and also ensures types
func (s SymbolTableVariableEntry) SymbolTableString() string {
	return fmt.Sprintf("%%%s", s.MemoryLocation)
}

// SymbolTable is an array of SymbolTableEntries
var SymbolTable []SymbolTableEntry

var symbolMap = make(map[string]SymbolTableEntry)

func InsertToSymbolMap(key string, value SymbolTableEntry) error {
	symbolMap[key] = value
	return nil
}

func AccSymbolMap(key string) (SymbolTableEntry, bool) {
	x, ok := symbolMap[key]
	return x, ok
}

// InsertToSymbolTable inserts a single entry into table
func InsertToSymbolTable(val string) (SymbolTableEntry, error) {
	if val, ok := symbolMap[val]; ok {
		return val, nil
	}
	var entry SymbolTableEntry
	switch val[0] {
	case '$':
		i, err := strconv.ParseInt(val[1:], 0, 32)
		if err != nil {
			return nil, err
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
		entry = &SymbolTableVariableEntry{
			MemoryLocation: val,
		}
	default:
		log.Fatalf("Unknown argument: %s", val)
	}
	SymbolTable = append(SymbolTable, entry)
	symbolMap[val] = entry
	return entry, nil
}

// GetRegs populates symbol table and gets virtual registers
func GetRegs(splitted []string, typ IRType, op IROp) (arg1, arg2, dst SymbolTableEntry, err error) {

	if typ == LBL {
		err = errors.New("we should never GetRegs for a label")
		return
	}

	if typ == BOP || typ == CBR || typ == LOP || typ == SOP || typ == DOP {
		if len(splitted) < 4 {
			err = errors.New("not enough args to a binary operand")
			return
		}
		dst, err = InsertToSymbolTable(splitted[1])
		if err != nil {
			return
		}
		arg1, err = InsertToSymbolTable(splitted[2])
		if err != nil {
			return
		}
		arg2, err = InsertToSymbolTable(splitted[3])
		if err != nil {
			return
		}
	} else if typ == UOP || typ == ASN {
		if len(splitted) < 3 {
			err = errors.New("not enough args to a unary operand")
			return
		}
		dst, err = InsertToSymbolTable(splitted[1])
		if err != nil {
			return
		}
		arg1, err = InsertToSymbolTable(splitted[2])
		if err != nil {
			return
		}
	} else if typ == JMP {
		if len(splitted) < 2 {
			err = errors.New("not enough args to a jump operand")
			return
		}
		arg1, err = InsertToSymbolTable(splitted[1])
		if err != nil {
			return
		}
	} else if typ == KEY {
		if !(op == RET || op == HALT) {
			if len(splitted) < 2 {
				err = errors.New("not enough args to a call/key operand")
				return
			}
			arg1, err = InsertToSymbolTable(splitted[1])
			if err != nil {
				return
			}
		}
	} else {
		err = fmt.Errorf("unknown instruction type: %d, %s", typ, op)
		return
	}
	return
}
