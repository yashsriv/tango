package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

// ModAssignment is used to evaluate an assignment statement
func ModAssignment(a, b, c Attrib) (*AddrCode, error) {
	op := string(b.(*token.Token).Lit)
	// TODO: Check el1 is addressable
	el1, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	if el1.Symbol == nil {
		return nil, fmt.Errorf("lhs must have a symbol table entry")
	}
	el2, ok := c.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", c)
	}
	if el2.Symbol == nil {
		return nil, fmt.Errorf("rhs must have a symbol table entry")
	}

	var irOp codegen.IROp
	var irType = codegen.BOP
	code := append(el1.Code, el2.Code...)
	switch op {
	case "|=":
		irOp = codegen.BOR
	case "+=":
		irOp = codegen.ADD
	case "-=":
		irOp = codegen.SUB
	case "^=":
		irOp = codegen.XOR
	case "*=":
		irOp = codegen.MUL
	case "/=":
		irType = codegen.DOP
		irOp = codegen.DIV
	case "%=":
		irType = codegen.DOP
		irOp = codegen.REM
	case "<<=":
		irType = codegen.SOP
		irOp = codegen.BSL
	case ">>=":
		irType = codegen.SOP
		irOp = codegen.BSR
	case "&=":
		irOp = codegen.BAND
	default:
		return nil, ErrUnsupported
	}
	code = append(code, codegen.IRIns{
		Typ:  irType,
		Op:   irOp,
		Dst:  el1.Symbol,
		Arg1: el1.Symbol,
		Arg2: el2.Symbol,
	})
	addrcode := &AddrCode{
		Code: code,
	}
	return addrcode, nil
}

// IncDec statement
func IncDec(a, b Attrib) (*AddrCode, error) {
	op := string(b.(*token.Token).Lit)
	// TODO: Check el1 is addressable
	el1, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	if el1.Symbol == nil {
		return nil, fmt.Errorf("lhs must have a symbol table entry")
	}
	code := append(el1.Code)
	var irOp codegen.IROp
	switch op {
	case "++":
		irOp = codegen.INC
	case "--":
		irOp = codegen.DEC
	}
	code = append(code, codegen.IRIns{
		Typ:  codegen.KEY,
		Op:   irOp,
		Arg1: el1.Symbol,
	})
	addrcode := &AddrCode{
		Code: code,
	}
	return addrcode, nil
}

// Decl is used to create a declaration statement
func Assignments(lhs, rhs Attrib) (*AddrCode, error) {

	// Obtain list of identifiers
	lhsList, ok := lhs.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", lhs)
	}

	// Obtain rhs of expression
	rhsList, ok := rhs.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", rhs)
	}

	// Number of elements in lhs and rhs must be same
	if len(lhsList) != len(rhsList) {
		return nil, fmt.Errorf("unequal number of elements in lhs and rhs: %d and %d", len(lhsList), len(rhsList))
	}

	// The code returned for this particular statement
	code := make([]codegen.IRIns, 0)

	// We assign values to temporary variables in case the rhs uses lhs values
	// This is not necessary if we only have one rhs and lhs as the single statement will
	// handle it gracefully
	entries := make([]codegen.SymbolTableEntry, len(rhsList))

	// If rhs is present and has more than one expression
	if len(rhsList) != 1 {
		// For each expression, store its value in a temporary variable
		for i, expr := range rhsList {
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
	for i, declName := range lhsList {
		// TODO: Check Addressable
		entry, ok := declName.Symbol.(*codegen.SymbolTableVariableEntry)
		if !ok {
			return nil, fmt.Errorf("lhs %s of expression should be a literal", declName)
		}
		entry.Declared = true

		// if there is a rhs
		entry.Assignments++
		var arg1 codegen.SymbolTableEntry
		// If number of expressions are greater than 1, they have
		// already been evaluated and assigned a value above.
		// Just refer to that value here.
		if len(rhsList) != 1 {
			arg1 = entries[i]
		} else {
			// Else evaluate expression and refer to its value here
			code = append(code, rhsList[i].Code...)
			arg1 = rhsList[i].Symbol
		}
		ins := codegen.IRIns{
			Typ:  codegen.ASN,
			Op:   codegen.ASNO,
			Dst:  declName.Symbol,
			Arg1: arg1,
		}
		code = append(code, ins)
	}

	addrcode := &AddrCode{
		Code: code,
	}

	return addrcode, nil
}
