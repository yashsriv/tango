package ast

import (
	"errors"
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
		return nil, fmt.Errorf("[ModAssignment] unable to type cast %v to *AddrCode", a)
	}
	if el1.Symbol == nil {
		return nil, fmt.Errorf("[ModAssignment] lhs must have a symbol table entry")
	}
	el2, ok := c.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[ModAssignment] unable to type cast %v to *AddrCode", c)
	}
	if el2.Symbol == nil {
		return nil, fmt.Errorf("[ModAssignment] rhs must have a symbol table entry")
	}

	var irOp codegen.IROp
	var irType = codegen.BOP
	code := append(el1.Code, el2.Code...)
	el1Evaluated := EvalWrapped(el1.Symbol)
	code = append(code, el1Evaluated.Code...)
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
	CheckOperandType(irOp, el1.Symbol.Type())
	CheckOperandType(irOp, el2.Symbol.Type())
	if !SameType(el1.Symbol.Type(), el2.Symbol.Type()) {
		return nil, errors.New("operands on either side of binary expression don't have the same type")
	}

	code = append(code, codegen.IRIns{
		Typ:  irType,
		Op:   irOp,
		Dst:  el1Evaluated.Symbol,
		Arg1: el1Evaluated.Symbol,
		Arg2: el2.Symbol,
	})
	if el1.Symbol.(*codegen.VariableEntry).Extra != nil {
		var ins codegen.IRIns
		switch e := el1.Symbol.(*codegen.VariableEntry).Extra.(type) {
		case ptrWrap:
			ins = codegen.IRIns{
				Typ:  codegen.KEY,
				Op:   codegen.PUT,
				Dst:  e.dest,
				Arg1: &codegen.LiteralEntry{Value: 0, LType: intType},
				Arg2: el1Evaluated.Symbol,
			}
		case arrWrap:
			ins = codegen.IRIns{
				Typ:  codegen.KEY,
				Op:   codegen.PUT,
				Dst:  e.dest,
				Arg1: e.index,
				Arg2: el1Evaluated.Symbol,
			}
		case structWrap:
			ins = codegen.IRIns{
				Typ:  codegen.KEY,
				Op:   codegen.PUT,
				Dst:  e.dest,
				Arg1: &codegen.LiteralEntry{Value: e.index, LType: intType},
				Arg2: el1Evaluated.Symbol,
			}
		default:
			return nil, ErrUnsupported
		}
		code = append(code, ins)
	}
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
		return nil, fmt.Errorf("[IncDec] unable to type cast %v to *AddrCode", a)
	}
	if el1.Symbol == nil {
		return nil, fmt.Errorf("[IncDec] lhs must have a symbol table entry")
	}
	code := append(el1.Code)
	var irOp codegen.IROp
	switch op {
	case "++":
		irOp = codegen.INC
	case "--":
		irOp = codegen.DEC
	}
	CheckOperandType(irOp, el1.Symbol.Type())
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
		return nil, fmt.Errorf("[Assignments] unable to typecast %v to []*AddrCode", lhs)
	}

	// Obtain rhs of expression
	rhsList, ok := rhs.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[Assignments] unable to typecast %v to []*AddrCode", rhs)
	}

	// Number of elements in lhs and rhs must be same
	if len(lhsList) != len(rhsList) {
		return nil, fmt.Errorf("[Assignments] unequal number of elements in lhs and rhs: %d and %d", len(lhsList), len(rhsList))
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

	for i := range lhsList {
		if !SameType(lhsList[i].Symbol.Type(), rhsList[i].Symbol.Type()) {
			return nil, fmt.Errorf("wrong type of rhs in assignment. expected %v, got %v", lhsList[i].Symbol.Type(), rhsList[i].Symbol.Type())
		}
	}

	// For each element in the rhs, perform some operations
	for i, declName := range lhsList {
		code = append(code, declName.Code...)
		// TODO: Check Addressable
		entry, ok := declName.Symbol.(*codegen.VariableEntry)
		if !ok {
			return nil, fmt.Errorf("[Assignments] lhs %v of expression should be a literal", declName)
		}

		if entry.Constant {
			return nil, fmt.Errorf("[Assignments] cannot assign to a constant")
		}

		// if there is a rhs
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
		var ins codegen.IRIns
		if entry.Extra != nil {
			switch e := entry.Extra.(type) {
			case ptrWrap:
				ins = codegen.IRIns{
					Typ:  codegen.KEY,
					Op:   codegen.PUT,
					Dst:  e.dest,
					Arg1: &codegen.LiteralEntry{Value: 0, LType: intType},
					Arg2: arg1,
				}
			case arrWrap:
				ins = codegen.IRIns{
					Typ:  codegen.KEY,
					Op:   codegen.PUT,
					Dst:  e.dest,
					Arg1: e.index,
					Arg2: arg1,
				}
			case structWrap:
				ins = codegen.IRIns{
					Typ:  codegen.KEY,
					Op:   codegen.PUT,
					Dst:  e.dest,
					Arg1: &codegen.LiteralEntry{Value: e.index, LType: intType},
					Arg2: arg1,
				}
			default:
				return nil, ErrUnsupported
			}
		} else {
			ins = codegen.IRIns{
				Typ:  codegen.ASN,
				Op:   codegen.ASNO,
				Dst:  declName.Symbol,
				Arg1: arg1,
			}
		}
		code = append(code, ins)
	}

	addrcode := &AddrCode{
		Code: code,
	}

	return addrcode, nil
}
