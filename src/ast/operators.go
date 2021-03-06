package ast

import (
	"errors"
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

type arrWrap struct {
	dest  codegen.SymbolTableEntry
	index codegen.SymbolTableEntry
}

type ptrWrap struct {
	dest codegen.SymbolTableEntry
}

type structWrap struct {
	dest  codegen.SymbolTableEntry
	index int
}

// EvalWrapped evaluates a wrapped up type
func EvalWrapped(el Attrib) *AddrCode {
	addrCode := el.(*AddrCode)
	symbol, ok := addrCode.Symbol.(*codegen.VariableEntry)
	if !ok || symbol.Extra == nil {
		return addrCode
	}

	code := addrCode.Code
	switch e := symbol.Extra.(type) {
	case ptrWrap:
		code = append(code, codegen.IRIns{
			Typ:  codegen.UOP,
			Op:   codegen.VAL,
			Dst:  symbol,
			Arg1: e.dest,
		})
	case arrWrap:
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.TAKE,
			Dst:  symbol,
			Arg1: e.dest,
			Arg2: e.index,
		})
	case structWrap:
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.TAKE,
			Dst:  symbol,
			Arg1: e.dest,
			Arg2: &codegen.LiteralEntry{Value: e.index},
		})
	default:
		panic(ErrUnsupported)
	}

	return &AddrCode{Symbol: symbol, Code: code}
}

// EvalStructAccess evaluates struct access
func EvalStructAccess(a, b Attrib) (*AddrCode, error) {
	el, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[EvalArrAccess] unable to type cast %v to *AddrCode", a)
	}

	identifier := string(b.(*token.Token).Lit)

	structType, isStruct := el.Symbol.Type().(codegen.StructType)
	if !isStruct {
		return nil, fmt.Errorf("access operation on non-struct type: %s", el.Symbol.Type())
	}

	index, ok := structType.FieldMap[identifier]
	if !ok {
		if structType.Name != "" {
			methodName := fmt.Sprintf("0%s_%s", structType, identifier)
			entry, err := codegen.SymbolTable.GetSymbol(methodName)
			if err == nil {
				if _, isTarget := entry.(*codegen.TargetEntry); isTarget {
					return &AddrCode{Symbol: entry, Code: append(el.Code,
						codegen.IRIns{
							Typ:  codegen.KEY,
							Op:   codegen.PARAM,
							Arg1: el.Symbol,
						},
					)}, nil
				}
			}
		}
		return nil, fmt.Errorf("unknown field %s in %v", identifier, structType)
	}

	code := el.Code
	entry := CreateTemporary(structType.FieldTypes[index])
	code = append(code, entry.Code...)

	if lvalMode {
		entry.Symbol.(*codegen.VariableEntry).Extra = structWrap{dest: el.Symbol, index: index}
	} else {
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.TAKE,
			Dst:  entry.Symbol,
			Arg1: el.Symbol,
			Arg2: &codegen.LiteralEntry{Value: index},
		})
	}

	return &AddrCode{Symbol: entry.Symbol, Code: code}, nil
}

// EvalArrAccess evaluates array access
func EvalArrAccess(a, b Attrib) (*AddrCode, error) {
	el, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[EvalArrAccess] unable to type cast %v to *AddrCode", a)
	}

	arrType, isArr := el.Symbol.Type().(codegen.ArrType)
	if !isArr {
		return nil, fmt.Errorf("indexing operation on non-array type: %s", el.Symbol.Type())
	}

	index, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[EvalArrAccess] unable to typecast %v to *AddrCode", b)
	}

	code := el.Code
	code = append(code, el.Code...)
	code = append(code, index.Code...)

	entry := CreateTemporary(arrType.Of)
	code = append(code, entry.Code...)

	if lvalMode {
		entry.Symbol.(*codegen.VariableEntry).Extra = arrWrap{dest: el.Symbol, index: index.Symbol}
	} else {
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.TAKE,
			Dst:  entry.Symbol,
			Arg1: el.Symbol,
			Arg2: index.Symbol,
		})
	}

	return &AddrCode{Symbol: entry.Symbol, Code: code}, nil
}

// EvalArrSlice evaluates array slicing
func EvalArrSlice(a, b, c Attrib) (*AddrCode, error) {
	return nil, ErrUnsupported
}

// PointerOp evaluates pointer operation
func PointerOp(op string, el *AddrCode) (*AddrCode, error) {
	var entry *AddrCode
	switch op {
	case "&":
		code := el.Code
		entry = CreateTemporary(codegen.PtrType{To: el.Symbol.Type()})
		code = append(code, entry.Code...)
		code = append(code, codegen.IRIns{
			Typ:  codegen.UOP,
			Op:   codegen.ADDR,
			Dst:  entry.Symbol,
			Arg1: el.Symbol,
		})
		addrcode := &AddrCode{
			Symbol: entry.Symbol,
			Code:   code,
		}
		return addrcode, nil
	case "*":
		ptr, isPtr := el.Symbol.Type().(codegen.PtrType)
		if !isPtr {
			return nil, fmt.Errorf("trying to dereference non-pointer type %v", el.Symbol.Type())
		}
		entry = CreateTemporary(ptr.To)
		code := el.Code
		code = append(code, entry.Code...)
		if lvalMode {
			entry.Symbol.(*codegen.VariableEntry).Extra = ptrWrap{dest: el.Symbol}
		} else {
			code = append(code, codegen.IRIns{
				Typ:  codegen.UOP,
				Op:   codegen.VAL,
				Dst:  entry.Symbol,
				Arg1: el.Symbol,
			})
		}
		addrcode := &AddrCode{
			Symbol: entry.Symbol,
			Code:   code,
		}
		return addrcode, nil
	default:
		return nil, ErrUnsupported
	}
}

// UnaryOp generates code for a unary expression
func UnaryOp(a Attrib, b Attrib) (*AddrCode, error) {
	op := string(a.(*token.Token).Lit)
	el, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}
	entry := CreateTemporary(el.Symbol.Type())
	code := el.Code
	code = append(code, entry.Code...)
	var irOp codegen.IROp
	switch op {
	case "*", "&":
		return PointerOp(op, el)
	case "+":
		irOp = codegen.ADD
	case "-":
		irOp = codegen.NEG
	case "^":
		irOp = codegen.BNOT
	case "!":
		irOp = codegen.NOT
	default:
		return nil, ErrUnsupported
	}
	CheckOperandType(irOp, el.Symbol.Type())
	if irOp == codegen.ADD {
		return el, nil
	}
	code = append(code, codegen.IRIns{
		Typ:  codegen.UOP,
		Op:   irOp,
		Dst:  entry.Symbol,
		Arg1: el.Symbol,
	})
	addrcode := &AddrCode{
		Symbol: entry.Symbol,
		Code:   code,
	}
	return addrcode, nil
}

// BinaryOp generates code for a binary expression
func BinaryOp(a Attrib, b Attrib, c Attrib) (*AddrCode, error) {
	op := string(b.(*token.Token).Lit)
	el1, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	el2, ok := c.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", c)
	}
	entry := CreateTemporary(el1.Symbol.Type())
	code := append(el1.Code, el2.Code...)
	code = append(code, entry.Code...)
	var irOp codegen.IROp
	var irType = codegen.BOP
	switch op {
	case "+":
		irOp = codegen.ADD
	case "-":
		irOp = codegen.SUB
	case "*":
		irOp = codegen.MUL
	case "&&":
		irOp = codegen.AND
	case "||":
		irOp = codegen.OR
	case "&":
		irOp = codegen.BAND
	case "|":
		irOp = codegen.BOR
	case "^":
		irOp = codegen.XOR
	case "/":
		irType = codegen.DOP
		irOp = codegen.DIV
	case "%":
		irType = codegen.DOP
		irOp = codegen.REM
	case "<<":
		irType = codegen.SOP
		irOp = codegen.BSL
	case ">>":
		irType = codegen.SOP
		irOp = codegen.BSR
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
		Dst:  entry.Symbol,
		Arg1: el1.Symbol,
		Arg2: el2.Symbol,
	})
	addrcode := &AddrCode{
		Symbol: entry.Symbol,
		Code:   code,
	}
	return addrcode, nil
}

var relOpCount = 0

// RelOp generates code for a binary expression
func RelOp(a Attrib, op string, c Attrib) (*AddrCode, error) {
	el1, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	el2, ok := c.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", c)
	}
	if _, isLiteral := el1.Symbol.(*codegen.LiteralEntry); isLiteral {
		tmp := CreateTemporary(el1.Symbol.Type())
		el1.Code = append(el1.Code, tmp.Code...)
		el1.Code = append(el1.Code, codegen.IRIns{
			Typ:  codegen.ASN,
			Op:   codegen.ASNO,
			Dst:  tmp.Symbol,
			Arg1: el1.Symbol,
		})
		el1.Symbol = tmp.Symbol
	}

	if _, isLiteral := el2.Symbol.(*codegen.LiteralEntry); isLiteral {
		tmp := CreateTemporary(el2.Symbol.Type())
		el2.Code = append(el2.Code, tmp.Code...)
		el2.Code = append(el2.Code, codegen.IRIns{
			Typ:  codegen.ASN,
			Op:   codegen.ASNO,
			Dst:  tmp.Symbol,
			Arg1: el2.Symbol,
		})
		el2.Symbol = tmp.Symbol
	}

	trueLbl := &codegen.TargetEntry{
		Target: fmt.Sprintf("_rel_op_%d_true", relOpCount),
	}

	endLbl := &codegen.TargetEntry{
		Target: fmt.Sprintf("_rel_op_%d_end", relOpCount),
	}

	relOpCount++

	entry := CreateTemporary(boolType)
	code := append(el1.Code, el2.Code...)
	code = append(code, entry.Code...)
	var irOp codegen.IROp
	switch op {
	case "<":
		irOp = codegen.BRLT
	case "<=":
		irOp = codegen.BRLTE
	case ">":
		irOp = codegen.BRGT
	case ">=":
		irOp = codegen.BRGTE
	case "==":
		irOp = codegen.BREQ
	case "!=":
		irOp = codegen.BRNEQ
	default:
		return nil, ErrUnsupported
	}

	CheckOperandType(irOp, el1.Symbol.Type())
	CheckOperandType(irOp, el2.Symbol.Type())
	if !SameType(el1.Symbol.Type(), el2.Symbol.Type()) {
		return nil, errors.New("operands on either side of binary expression don't have the same type")
	}

	// Comparison and jmp to true label
	code = append(code, codegen.IRIns{
		Typ:  codegen.CBR,
		Op:   irOp,
		Dst:  trueLbl,
		Arg1: el2.Symbol,
		Arg2: el1.Symbol,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.ASN,
		Op:  codegen.ASNO,
		Dst: entry.Symbol,
		Arg1: &codegen.LiteralEntry{
			Value: 0,
			LType: boolType,
		},
	})
	code = append(code, codegen.IRIns{
		Typ:  codegen.JMP,
		Op:   codegen.JMPO,
		Arg1: endLbl,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: trueLbl,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.ASN,
		Op:  codegen.ASNO,
		Dst: entry.Symbol,
		Arg1: &codegen.LiteralEntry{
			Value: 1,
			LType: boolType,
		},
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: endLbl,
	})
	addrcode := &AddrCode{
		Symbol: entry.Symbol,
		Code:   code,
	}
	return addrcode, nil
}

func AndOp(a, b Attrib) (*AddrCode, error) {
	el1, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	el2, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}
	falseLbl := &codegen.TargetEntry{
		Target: fmt.Sprintf("_rel_op_%d_false", relOpCount),
	}

	endLbl := &codegen.TargetEntry{
		Target: fmt.Sprintf("_rel_op_%d_end", relOpCount),
	}
	CheckOperandType(codegen.AND, el1.Symbol.Type())
	CheckOperandType(codegen.AND, el2.Symbol.Type())
	if !SameType(el1.Symbol.Type(), el2.Symbol.Type()) {
		return nil, errors.New("operands on either side of binary expression don't have the same type")
	}

	relOpCount++

	entry := CreateTemporary(boolType)
	code := append(el1.Code)
	code = append(code, entry.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BRNEQ,
		Dst: falseLbl,
		Arg1: &codegen.LiteralEntry{
			Value: 1,
			LType: boolType,
		},
		Arg2: el1.Symbol,
	})
	code = append(code, el2.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BRNEQ,
		Dst: falseLbl,
		Arg1: &codegen.LiteralEntry{
			Value: 1,
			LType: boolType,
		},
		Arg2: el2.Symbol,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.ASN,
		Op:  codegen.ASNO,
		Dst: entry.Symbol,
		Arg1: &codegen.LiteralEntry{
			Value: 1,
			LType: boolType,
		},
	})
	code = append(code, codegen.IRIns{
		Typ:  codegen.JMP,
		Op:   codegen.JMPO,
		Arg1: endLbl,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: falseLbl,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.ASN,
		Op:  codegen.ASNO,
		Dst: entry.Symbol,
		Arg1: &codegen.LiteralEntry{
			Value: 0,
			LType: boolType,
		},
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: endLbl,
	})
	addrcode := &AddrCode{
		Symbol: entry.Symbol,
		Code:   code,
	}
	return addrcode, nil
}

func OrOp(a, b Attrib) (*AddrCode, error) {
	el1, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	el2, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}
	trueLbl := &codegen.TargetEntry{
		Target: fmt.Sprintf("_rel_op_%d_true", relOpCount),
	}
	endLbl := &codegen.TargetEntry{
		Target: fmt.Sprintf("_rel_op_%d_end", relOpCount),
	}

	CheckOperandType(codegen.OR, el1.Symbol.Type())
	CheckOperandType(codegen.OR, el2.Symbol.Type())
	if !SameType(el1.Symbol.Type(), el2.Symbol.Type()) {
		return nil, errors.New("operands on either side of binary expression don't have the same type")
	}
	relOpCount++

	entry := CreateTemporary(boolType)
	code := append(el1.Code)
	code = append(code, entry.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BREQ,
		Dst: trueLbl,
		Arg1: &codegen.LiteralEntry{
			Value: 1,
			LType: boolType,
		},
		Arg2: el1.Symbol,
	})
	code = append(code, el2.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BREQ,
		Dst: trueLbl,
		Arg1: &codegen.LiteralEntry{
			Value: 1,
			LType: boolType,
		},
		Arg2: el2.Symbol,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.ASN,
		Op:  codegen.ASNO,
		Dst: entry.Symbol,
		Arg1: &codegen.LiteralEntry{
			Value: 0,
			LType: boolType,
		},
	})
	code = append(code, codegen.IRIns{
		Typ:  codegen.JMP,
		Op:   codegen.JMPO,
		Arg1: endLbl,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: trueLbl,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.ASN,
		Op:  codegen.ASNO,
		Dst: entry.Symbol,
		Arg1: &codegen.LiteralEntry{
			Value: 1,
			LType: boolType,
		},
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: endLbl,
	})
	addrcode := &AddrCode{
		Symbol: entry.Symbol,
		Code:   code,
	}
	return addrcode, nil
}

// ProcessName is used to process a name
func ProcessName(a Attrib) (*AddrCode, error) {
	switch v := a.(type) {
	case *codegen.VariableEntry:
		return &AddrCode{Symbol: v}, nil
	case *codegen.TargetEntry:
		return &AddrCode{Symbol: v}, nil
	default:
		fmt.Printf("%T\n", v)
		return nil, ErrShouldBeVariable
	}
}
