package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

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
	var ncode codegen.IRIns
	switch op {
	case "-":
		ncode = codegen.IRIns{
			Typ:  codegen.UOP,
			Op:   codegen.NEG,
			Dst:  entry,
			Arg1: el.Symbol,
		}
	case "^":
		ncode = codegen.IRIns{
			Typ:  codegen.UOP,
			Op:   codegen.BNOT,
			Dst:  entry,
			Arg1: el.Symbol,
		}
	default:
		return nil, ErrUnsupported
	}
	code = append(code, ncode)
	addrcode := &AddrCode{
		Symbol: entry,
		Code:   code,
	}
	return addrcode, nil
}

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
	var ncode codegen.IRIns
	switch op {
	case "+":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.ADD,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case "-":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.SUB,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case "*":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.MUL,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case "^":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.XOR,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case "|":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.BOR,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case "/":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.DIV,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case "%":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.REM,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case "<<":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.BSL,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case ">>":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.BSR,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	case "&":
		ncode = codegen.IRIns{
			Typ:  codegen.BOP,
			Op:   codegen.BAND,
			Dst:  entry,
			Arg1: el1.Symbol,
			Arg2: el2.Symbol,
		}
	default:
		return nil, ErrUnsupported
	}
	code = append(code, ncode)
	addrcode := &AddrCode{
		Symbol: entry,
		Code:   code,
	}
	return addrcode, nil
}
