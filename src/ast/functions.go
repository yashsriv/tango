package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

type ArgType struct {
	ArgName string
}

// FuncSign is the function's signature
func FuncSign(a, b, c, d Attrib) (*AddrCode, error) {
	// TODO: Handle other stuff like arg list, return type and method declarations
	identifier := string(a.(*token.Token).Lit)
	start := &codegen.TargetEntry{
		Target: fmt.Sprintf("_func_%s", identifier),
	}
	// Associating identifier with some entry
	err := codegen.SymbolTable.InsertSymbol(identifier, start)
	if err != nil {
		return nil, err
	}

	NewScope()

	args := b.([]*ArgType)

	for i, arg := range args {
		symbol := &codegen.VariableEntry{
			MemoryLocation: codegen.StackMemory{BaseOffset: 4 * (i + 2)},
			Name:           arg.ArgName,
		}
		err = codegen.SymbolTable.InsertSymbol(arg.ArgName, symbol)
		if err != nil {
			return nil, err
		}
		codegen.CreateAddrDescEntry(symbol)
	}

	return &AddrCode{Symbol: start}, err
}

// EvalArgType evaluates an argument to function
func EvalArgType(a, b Attrib) (*ArgType, error) {
	identifier := string(a.(*token.Token).Lit)
	return &ArgType{ArgName: identifier}, nil
}
