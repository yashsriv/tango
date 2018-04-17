package codegen

import "fmt"

// SymbolTableEntry is an entry in the SymbolTable
type SymbolTableEntry interface {
	symbolTableEntryDummy()
	String() string
}

// LiteralEntry refers to a literal in the symbol table
type LiteralEntry struct {
	Value int
}

func (*LiteralEntry) symbolTableEntryDummy() {}

func (l *LiteralEntry) String() string {
	return fmt.Sprintf("$%d", l.Value)
}

// TargetEntry refers to a target in the symbol table
type TargetEntry struct {
	Target string
}

func (*TargetEntry) symbolTableEntryDummy() {}

func (t *TargetEntry) String() string {
	return fmt.Sprintf("#%s", t.Target)
}

// VariableEntry refers to a register in the symbol table
type VariableEntry struct {
	Constant       bool
	MemoryLocation MemoryLocation
	Name           string
}

func (*VariableEntry) symbolTableEntryDummy() {}

func (v *VariableEntry) String() string {
	return fmt.Sprintf("%s", v.Name)
}
