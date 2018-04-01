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
