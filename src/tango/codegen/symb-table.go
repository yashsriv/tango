package codegen

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
