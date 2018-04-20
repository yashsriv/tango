package codegen

import "fmt"

// SymbolTableEntry is an entry in the SymbolTable
type SymbolTableEntry interface {
	symbolTableEntryDummy()
	Type() TypeEntry
	String() string
}

// LiteralEntry refers to a literal in the symbol table
type LiteralEntry struct {
	Value int
	LType TypeEntry
}

func (*LiteralEntry) symbolTableEntryDummy() {}
func (l *LiteralEntry) String() string {
	return fmt.Sprintf("$%d", l.Value)
}

// Type is necessary
func (l *LiteralEntry) Type() TypeEntry {
	return l.LType
}

// TargetEntry refers to a target in the symbol table
type TargetEntry struct {
	Target  string
	RetType TypeEntry
	InType  []TypeEntry
}

func (*TargetEntry) symbolTableEntryDummy() {}
func (t *TargetEntry) String() string {
	return fmt.Sprintf("#%s", t.Target)
}

// Type is necessary
func (t *TargetEntry) Type() TypeEntry {
	return t.RetType
}

// VariableEntry refers to a register in the symbol table
type VariableEntry struct {
	Constant       bool
	MemoryLocation MemoryLocation
	Name           string
	VType          TypeEntry
	Extra          interface{}
}

func (*VariableEntry) symbolTableEntryDummy() {}
func (v *VariableEntry) String() string {
	return fmt.Sprintf("%s", v.Name)
}

// Type is necessary
func (v *VariableEntry) Type() TypeEntry {
	return v.VType
}

// TypeEntry represents a type in the symbol table
type TypeEntry interface {
	SymbolTableEntry
	typeEntryDummy()
}

// BasicType represents an inbuilt word-aligned type
type BasicType struct {
	Name string
}

func (*BasicType) symbolTableEntryDummy() {}
func (*BasicType) typeEntryDummy()        {}
func (b *BasicType) String() string {
	return fmt.Sprintf("%s", b.Name)
}

// Type is necessary
func (b *BasicType) Type() TypeEntry {
	return b
}

// VoidType represents a void type
type VoidType struct {
}

func (*VoidType) symbolTableEntryDummy() {}
func (*VoidType) typeEntryDummy()        {}
func (*VoidType) String() string {
	return "VoidType"
}

// Type is necessary
func (v *VoidType) Type() TypeEntry {
	return v
}

// PtrType is a pointer to some other type
type PtrType struct {
	To TypeEntry
}

func (PtrType) symbolTableEntryDummy() {}
func (PtrType) typeEntryDummy()        {}
func (p PtrType) String() string {
	return fmt.Sprintf("pointer to %s", p.To)
}

// Type is necessary
func (p PtrType) Type() TypeEntry {
	return p
}

// ArrType is an array of other types
type ArrType struct {
	Of   TypeEntry
	Size int
}

func (ArrType) symbolTableEntryDummy() {}
func (ArrType) typeEntryDummy()        {}
func (a ArrType) String() string {
	return fmt.Sprintf("array of %s", a.Of)
}

// Type is necessary
func (a ArrType) Type() TypeEntry {
	return a
}
