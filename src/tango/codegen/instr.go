package codegen

// IRLabel represents a jumpable location
type IRLabel string

// IRIns represents an IR Instruction
type IRIns struct {
	Typ   IRType
	Op    IROp
	Arg1  SymbolTableEntry
	Arg2  SymbolTableEntry
	Dst   SymbolTableEntry
	Label IRLabel
}
