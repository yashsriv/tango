package codegen

import (
	"errors"
	"fmt"
)

// SymbolTableEntry is an entry in the SymbolTable
type SymbolTableEntry interface {
	SymbolTableString() string
}

// SymbolTableLiteralEntry refers to a literal in the symbol table
type SymbolTableLiteralEntry struct {
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

// symbolTable represents a table
type symbolTable struct {
	symbolMap map[string]SymbolTableEntry
	parent    *symbolTable
}

// SymbolTable is current table
var SymbolTable *symbolTable

// rootTable refers to the global rootTable
var rootTable *symbolTable

// tableStack maintains a stack of symbol tables
var tableStack []*symbolTable

// Initialize data structures
func init() {
	rootTable = &symbolTable{
		symbolMap: make(map[string]SymbolTableEntry),
		parent:    nil,
	}
	SymbolTable = rootTable
}

// push to table stack
func pushToStack() {
	tableStack = append(tableStack, SymbolTable)
}

// pop from table stack
func popFromStack() (table *symbolTable, err error) {
	l := len(tableStack)
	if l == 0 {
		err = ErrEmptyTableStack
		return
	}
	table = tableStack[l-1]
	tableStack = tableStack[:l-1]
	return
}

// ErrAlreadyExists is error when symbol already exists in table
var ErrAlreadyExists = errors.New("symbol already exists in table")

// ErrDoesntExist is error when symbol is not in table
var ErrDoesntExist = errors.New("symbol doesn't exist in table")

// ErrEmptyTableStack is an error thrown when trying to pop off empty table stack
var ErrEmptyTableStack = errors.New("expected tableStack to never be empty")

func (s *symbolTable) InsertSymbol(key string, value SymbolTableEntry) error {
	// Check if already exists in current scope
	_, ok := s.symbolMap[key]
	if ok {
		return ErrAlreadyExists
	}

	// Insert otherwise
	s.symbolMap[key] = value
	return nil
}

func (s *symbolTable) GetSymbol(key string) (SymbolTableEntry, error) {
	// Check if symbol exists in current scope
	x, ok := s.symbolMap[key]
	if !ok {
		if s.parent != nil {
			// If not, check in higher scopes
			return s.parent.GetSymbol(key)
		}
		return nil, ErrDoesntExist
	}
	// If highest scope, then result of this layer is the result
	return x, nil
}

// NewScope creates a new scope
func NewScope() {
	pushToStack()
	SymbolTable = &symbolTable{
		symbolMap: make(map[string]SymbolTableEntry),
		parent:    SymbolTable,
	}
}

// EndScope ends a scope
func EndScope() (err error) {
	SymbolTable, err = popFromStack()
	return
}
