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
	var methodArgType *ArgType
	var methodName string
	if d != nil {
		argl := d.([]*ArgType)
		if len(argl) != 1 {
			return nil, fmt.Errorf("only 1 param allowed in method. got %d", len(argl))
		}
		methodArgType = argl[0]
		t, isStruct := methodArgType.Type.(codegen.StructType)
		if !isStruct {
			return nil, fmt.Errorf("methods can only be defined on structs. got %s", methodArgType.Type)
		}
		if _, ok := t.FieldMap[identifier]; ok {
			return nil, fmt.Errorf("redeclaring field %s of %v", identifier, t)
		}
		if t.Name == "" {
			return nil, fmt.Errorf("can only declare methods on named structs")
		}
		methodName = t.Name
	}
	var start *codegen.TargetEntry
	if methodArgType != nil {
		start = &codegen.TargetEntry{
			Target:  fmt.Sprintf("_func_0%s_%s", methodName, identifier),
			RetType: retType,
			InType:  inTypes,
		}
	} else {
		start = &codegen.TargetEntry{
			Target:  fmt.Sprintf("_func_%s", identifier),
			RetType: retType,
			InType:  inTypes,
		}
	}
	// Associating identifier with some entry
	var err error
	if methodArgType == nil {
		err = codegen.SymbolTable.InsertSymbol(identifier, start)
	} else {
		err = codegen.SymbolTable.InsertSymbol("0"+methodArgType.Type.String()+"_"+identifier, start)
	}
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

	if methodArgType != nil {
		symbol := &codegen.VariableEntry{
			MemoryLocation: codegen.StackMemory{BaseOffset: 4 * (len(args) + 2)},
			Name:           methodArgType.ArgName,
			VType:          methodArgType.Type,
		}
		err = codegen.SymbolTable.InsertSymbol(methodArgType.ArgName, symbol)
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

func FuncType(a, b Attrib) (Attrib, error) {
	return nil, nil
}
