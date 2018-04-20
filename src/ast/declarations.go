package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

// Decl is used to create a declaration statement
func Decl(declnamelist, types, exprlist Attrib, isConst bool) (*AddrCode, error) {

	// Obtain list of identifiers
	declnamelistAs, ok := declnamelist.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", declnamelist)
	}

	// Obtain rhs of expression
	var exprlistAs []*AddrCode

	// exprlist can be nil in case of var declaration so we must first check
	if exprlist != nil {
		exprlistAs, ok = exprlist.([]*AddrCode)
		if !ok {
			return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", exprlist)
		}

		// Number of elements in lhs and rhs must be same
		if len(declnamelistAs) != len(exprlistAs) {
			return nil, fmt.Errorf("unequal number of elements in lhs and rhs: %d and %d", len(declnamelistAs), len(exprlistAs))
		}

		for i := range exprlistAs {
			if !SameType(types.(codegen.TypeEntry), exprlistAs[i].Symbol.Type()) {
				return nil, fmt.Errorf("wrong type of rhs in declaration. expected %v, got %v", types.(codegen.TypeEntry), exprlistAs[i].Symbol.Type())
			}
		}
	} else if isConst {
		// for a constant declaration, rhs is compulsory
		return nil, fmt.Errorf("constant declaration must have rhs")
	}

	// The code returned for this particular statement
	code := make([]codegen.IRIns, 0)

	// We assign values to temporary variables in case the rhs uses lhs values
	// This is not necessary if we only have one rhs and lhs as the single statement will
	// handle it gracefully
	entries := make([]codegen.SymbolTableEntry, len(exprlistAs))

	// If rhs is present and has more than one expression
	if exprlist != nil && len(exprlistAs) != 1 {
		// For each expression, store its value in a temporary variable
		for i, expr := range exprlistAs {
			code = append(code, expr.Code...)
			entry := CreateTemporary(expr.Symbol.Type())
			entries[i] = entry.Symbol
			code = append(code, entry.Code...)
			ins := codegen.IRIns{
				Typ:  codegen.ASN,
				Op:   codegen.ASNO,
				Dst:  entry.Symbol,
				Arg1: expr.Symbol,
			}
			code = append(code, ins)
		}
	}

	// For each element in the rhs, perform some operations
	for i, declName := range declnamelistAs {
		declName.Symbol.(*codegen.VariableEntry).Constant = isConst
		declName.Symbol.(*codegen.VariableEntry).VType = types.(codegen.TypeEntry)
		code = append(code, declName.Code...)
		// if there is a rhs
		if exprlist != nil {
			var arg1 codegen.SymbolTableEntry
			// If number of expressions are greater than 1, they have
			// already been evaluated and assigned a value above.
			// Just refer to that value here.
			if len(exprlistAs) != 1 {
				arg1 = entries[i]
			} else {
				// Else evaluate expression and refer to its value here
				code = append(code, exprlistAs[i].Code...)
				arg1 = exprlistAs[i].Symbol
			}
			ins := codegen.IRIns{
				Typ:  codegen.ASN,
				Op:   codegen.ASNO,
				Dst:  declName.Symbol,
				Arg1: arg1,
			}
			code = append(code, ins)
		} else {
			code = append(code, Alloc(declName.Symbol.(*codegen.VariableEntry))...)
		}
	}

	addrcode := &AddrCode{
		Code: code,
	}

	return addrcode, nil
}

// MultConstDecl is used for a block of multiple constant declarations
func MultConstDecl(decl, decllist Attrib) (*AddrCode, error) {
	mergedList, err := MergeCodeList(decllist)
	if err != nil {
		return nil, err
	}
	declAsAddr, ok := decl.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[MultConstDecl] unable to type cast %v to *AddrCode", decl)
	}
	code := append(declAsAddr.Code, mergedList.Code...)
	addrcode := &AddrCode{
		Code: code,
	}

	return addrcode, nil
}

// FuncDecl is a declaration of a function
func FuncDecl(a, b Attrib) (*AddrCode, error) {
	name, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[FuncDecl] unable to type cast %v to *AddrCode", a)
	}
	body, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[FuncDecl] unable to type cast %v to *AddrCode", b)
	}
	code := make([]codegen.IRIns, 1, len(body.Code)+1)
	code[0] = codegen.IRIns{
		Typ: codegen.LBL,
		Dst: name.Symbol,
	}
	code = append(code, body.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.KEY,
		Op:  codegen.RET,
	})
	addrcode := &AddrCode{
		Code: code,
	}
	currentRetType = nil
	return addrcode, nil
}

// EvalTypeDecl defines a new type
func EvalTypeDecl(a, b Attrib) (*AddrCode, error) {
	identifier := string(a.(*token.Token).Lit)
	symbol := b.(codegen.TypeEntry)
	if s, isStruct := symbol.(codegen.StructType); isStruct {
		s.Name = identifier
		symbol = s
	}
	err := codegen.SymbolTable.InsertSymbol(identifier, symbol)
	if err != nil {
		return nil, err
	}
	return &AddrCode{}, nil
}

// NewName creates a new symbol table entry for a variable
func NewName(a Attrib) (*AddrCode, error) {
	identifier := string(a.(*token.Token).Lit)
	var code = make([]codegen.IRIns, 0)
	var location codegen.MemoryLocation
	if codegen.SymbolTable.IsRoot() {
		location = codegen.GlobalMemory{Location: "v" + identifier}
	} else {
		// TODO: This should be from types or something else
		offset := codegen.SymbolTable.Alloc(4)
		location = codegen.StackMemory{BaseOffset: offset}
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.ALLOC,
			Arg1: &codegen.LiteralEntry{Value: 4, LType: intType},
		})
	}
	symbol := &codegen.VariableEntry{
		MemoryLocation: location,
		Name:           identifier,
	}
	err := codegen.SymbolTable.InsertSymbol(identifier, symbol)
	if err != nil {
		return nil, err
	}
	// Create AddrDesc Entry
	codegen.CreateAddrDescEntry(symbol)

	return &AddrCode{Symbol: symbol, Code: code}, nil
}

// Name gets table entry for a symbol
func Name(a Attrib) (symbol codegen.SymbolTableEntry, err error) {
	identifier := string(a.(*token.Token).Lit)
	symbol, err = codegen.SymbolTable.GetSymbol(identifier)
	if identifier == "s" {
		// fmt.Printf("found %s with type %T and location %s\n", identifier, symbol, symbol.(*codegen.VariableEntry).MemoryLocation)
	}
	return
}

// CreateTemporary creates a temporary variable
func CreateTemporary(types codegen.TypeEntry) *AddrCode {
	var location codegen.MemoryLocation
	var code = make([]codegen.IRIns, 0)
	if codegen.SymbolTable.IsRoot() {
		location = codegen.GlobalMemory{Location: fmt.Sprintf("rtmp%d", tempCount)}
	} else {
		// TODO: This should be from types or something else
		offset := codegen.SymbolTable.Alloc(4)
		location = codegen.StackMemory{BaseOffset: offset}
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.ALLOC,
			Arg1: &codegen.LiteralEntry{Value: 4, LType: intType},
		})
	}
	symbol := &codegen.VariableEntry{
		MemoryLocation: location,
		Name:           fmt.Sprintf("rtmp%d", tempCount),
		VType:          types,
	}
	codegen.CreateAddrDescEntry(symbol)
	tempCount++
	return &AddrCode{Symbol: symbol, Code: code}
}
