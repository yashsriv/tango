package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

type ArgType struct {
	ArgName string
	Type    codegen.TypeEntry
}

// FuncSign is the function's signature
func FuncSign(a, b, c, d Attrib) (*AddrCode, error) {
	// TODO: Handle other stuff like arg list, return type and method declarations
	identifier := string(a.(*token.Token).Lit)
	args := b.([]*ArgType)
	retType := c.(codegen.TypeEntry)
	inTypes := make([]codegen.TypeEntry, len(args))
	for i, arg := range args {
		inTypes[i] = arg.Type
	}
	start := &codegen.TargetEntry{
		Target:  fmt.Sprintf("_func_%s", identifier),
		RetType: retType,
		InType:  inTypes,
	}
	// Associating identifier with some entry
	err := codegen.SymbolTable.InsertSymbol(identifier, start)
	if err != nil {
		return nil, err
	}

	currentRetType = retType

	NewScope()

	for i, arg := range args {
		symbol := &codegen.VariableEntry{
			MemoryLocation: codegen.StackMemory{BaseOffset: 4 * (i + 2)},
			Name:           arg.ArgName,
			VType:          arg.Type,
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
	retType := b.(codegen.TypeEntry)
	return &ArgType{ArgName: identifier, Type: retType}, nil
}
