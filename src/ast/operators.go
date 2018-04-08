package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

// UnaryOp generates code for a unary expression
func UnaryOp(a Attrib, b Attrib) (*AddrCode, error) {
	op := string(a.(*token.Token).Lit)
	el, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}
	entry, err := codegen.InsertToSymbolTable(fmt.Sprintf("rtmp%d", tempCount))
	tempCount++
	if err != nil {
		return nil, err
	}
	code := el.Code
	var irOp codegen.IROp
	switch op {
	case "-":
		irOp = codegen.NEG
	case "^":
		irOp = codegen.BNOT
	default:
		return nil, ErrUnsupported
	}
	code = append(code, codegen.IRIns{
		Typ:  codegen.UOP,
		Op:   irOp,
		Dst:  entry,
		Arg1: el.Symbol,
	})
	addrcode := &AddrCode{
		Symbol: entry,
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
	entry, err := codegen.InsertToSymbolTable(fmt.Sprintf("rtmp%d", tempCount))
	tempCount++
	if err != nil {
		return nil, err
	}
	code := append(el1.Code, el2.Code...)
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
	code = append(code, codegen.IRIns{
		Typ:  irType,
		Op:   irOp,
		Dst:  entry,
		Arg1: el1.Symbol,
		Arg2: el2.Symbol,
	})
	addrcode := &AddrCode{
		Symbol: entry,
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

	trueLbl, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_rel_op_%d_true", relOpCount))
	if err != nil {
		return nil, err
	}

	endLbl, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_rel_op_%d_end", relOpCount))
	if err != nil {
		return nil, err
	}

	relOpCount++

	entry, err := codegen.InsertToSymbolTable(fmt.Sprintf("rtmp%d", tempCount))
	if err != nil {
		return nil, err
	}
	tempCount++
	code := append(el1.Code, el2.Code...)
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
		Dst: entry,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$0",
			Value: 0,
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
		Dst: entry,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$1",
			Value: 1,
		},
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: endLbl,
	})
	addrcode := &AddrCode{
		Symbol: entry,
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
	falseLbl, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_rel_op_%d_false", relOpCount))
	if err != nil {
		return nil, err
	}

	endLbl, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_rel_op_%d_end", relOpCount))
	if err != nil {
		return nil, err
	}

	relOpCount++

	entry, err := codegen.InsertToSymbolTable(fmt.Sprintf("rtmp%d", tempCount))
	if err != nil {
		return nil, err
	}
	tempCount++
	code := append(el1.Code)
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BRNEQ,
		Dst: falseLbl,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$1",
			Value: 1,
		},
		Arg2: el1.Symbol,
	})
	code = append(code, el2.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BRNEQ,
		Dst: falseLbl,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$1",
			Value: 1,
		},
		Arg2: el2.Symbol,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.ASN,
		Op:  codegen.ASNO,
		Dst: entry,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$1",
			Value: 1,
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
		Dst: entry,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$0",
			Value: 0,
		},
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: endLbl,
	})
	addrcode := &AddrCode{
		Symbol: entry,
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
	trueLbl, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_rel_op_%d_true", relOpCount))
	if err != nil {
		return nil, err
	}

	endLbl, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_rel_op_%d_end", relOpCount))
	if err != nil {
		return nil, err
	}

	relOpCount++

	entry, err := codegen.InsertToSymbolTable(fmt.Sprintf("rtmp%d", tempCount))
	if err != nil {
		return nil, err
	}
	tempCount++
	code := append(el1.Code)
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BREQ,
		Dst: trueLbl,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$1",
			Value: 1,
		},
		Arg2: el1.Symbol,
	})
	code = append(code, el2.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BREQ,
		Dst: trueLbl,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$1",
			Value: 1,
		},
		Arg2: el2.Symbol,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.ASN,
		Op:  codegen.ASNO,
		Dst: entry,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$0",
			Value: 0,
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
		Dst: entry,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Repr:  "$1",
			Value: 1,
		},
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: endLbl,
	})
	addrcode := &AddrCode{
		Symbol: entry,
		Code:   code,
	}
	return addrcode, nil
}

func ProcessName(a Attrib) (*AddrCode, error) {
	identifier := string(a.(*token.Token).Lit)
	entry, ok := codegen.AccSymbolMap(identifier)
	if !ok {
		return nil, fmt.Errorf("Identifier %s not declared yet", identifier)
	}
	return &AddrCode{
		Symbol: entry,
	}, nil
}
