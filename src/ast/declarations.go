package ast

import (
	"fmt"
	"tango/src/codegen"
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

	// exprlist can be nil in case of var declration so we must first check
	if exprlist != nil {
		exprlistAs, ok = exprlist.([]*AddrCode)
		if !ok {
			return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", exprlist)
		}

		// Number of elements in lhs and rhs must be same
		if len(declnamelistAs) != len(exprlistAs) {
			return nil, fmt.Errorf("unequal number of elements in lhs and rhs: %d and %d", len(declnamelistAs), len(exprlistAs))
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
			entry, err := codegen.InsertToSymbolTable(fmt.Sprintf("rtmp%d", tempCount))
			entries[i] = entry
			tempCount++
			if err != nil {
				return nil, err
			}
			ins := codegen.IRIns{
				Typ:  codegen.ASN,
				Op:   codegen.ASNO,
				Dst:  entry,
				Arg1: expr.Symbol,
			}
			code = append(code, ins)
		}
	}

	// For each element in the rhs, perform some operations
	for i, declName := range declnamelistAs {
		entry, ok := declName.Symbol.(*codegen.SymbolTableVariableEntry)
		if !ok {
			return nil, fmt.Errorf("lhs %s of expression should be a literal", declName)
		}
		if entry.Declared {
			return nil, fmt.Errorf("%s is being declared twice in this scope", entry.SymbolTableString())
		}
		entry.Declared = true

		// if there is a rhs
		if exprlist != nil {
			entry.Assignments++
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
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", decl)
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
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	body, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
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
	return addrcode, nil
}
