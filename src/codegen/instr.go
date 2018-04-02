package codegen

import "fmt"

// IRLabel represents a jumpable location
type IRLabel string

// IRIns represents an IR Instruction
type IRIns struct {
	Typ  IRType
	Op   IROp
	Arg1 SymbolTableEntry
	Arg2 SymbolTableEntry
	Dst  SymbolTableEntry
}

func (i IRIns) String() string {
	var arg1 = ""
	if i.Arg1 != nil {
		arg1 = i.Arg1.SymbolTableString()
	}
	var arg2 = ""
	if i.Arg2 != nil {
		arg2 = i.Arg2.SymbolTableString()
	}
	var dst = ""
	if i.Dst != nil {
		dst = i.Dst.SymbolTableString()
	}
	if i.Typ == LBL {
		return fmt.Sprintf("%s:", i.Dst.(*SymbolTableTargetEntry).Target)
	}
	return fmt.Sprintf("%s,%s,%s,%s", i.Op, dst, arg1, arg2)
}
