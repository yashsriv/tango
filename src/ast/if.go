package ast

import (
	"fmt"
	"tango/src/codegen"
)

var ifElseCount = 0

// EvalIfHeader evaluates the IfHeader
func EvalIfHeader(a, b Attrib) (*AddrCode, error) {
	stmt, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	expr, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}
	code := append(stmt.Code, expr.Code...)
	addrCode := &AddrCode{
		Code:   code,
		Symbol: expr.Symbol,
	}
	return addrCode, nil
}

type ifElse struct {
	expr *AddrCode
	body *AddrCode
}

// EvalElseIf evaluates else if
func EvalElseIf(a, b Attrib) (*ifElse, error) {
	expr, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	body, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}
	return &ifElse{expr, body}, nil
}

// EvalIf evaluates if
func EvalIf(a, b, c, d Attrib) (*AddrCode, error) {

	ifexpr, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	ifbody, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}

	elseifList, ok := c.([]*ifElse)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to []*ifElse", c)
	}
	elseBody, ok := d.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", d)
	}

	end, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_if_else_%d_end", ifElseCount))
	if err != nil {
		return nil, err
	}
	entry, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_if_%d_end", ifElseCount))
	if err != nil {
		return nil, err
	}
	code := make([]codegen.IRIns, 0)

	// Evaluate Expression
	code = append(code, ifexpr.Code...)

	// Check Expression
	code = append(code, codegen.IRIns{
		Typ: codegen.CBR,
		Op:  codegen.BRNEQ,
		Arg1: &codegen.SymbolTableLiteralEntry{
			Value: 1,
			Repr:  "$1",
		},
		Arg2: ifexpr.Symbol,
		Dst:  entry,
	})

	// If Body
	code = append(code, ifbody.Code...)
	code = append(code, codegen.IRIns{
		Typ:  codegen.JMP,
		Op:   codegen.JMPO,
		Arg1: end,
	})
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: entry,
	})

	for i, v := range elseifList {
		entry, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_if_%d_else_if_%d_end", ifElseCount, i))
		if err != nil {
			return nil, err
		}
		// Evaluate Expression
		code = append(code, v.expr.Code...)

		// Check Expression
		code = append(code, codegen.IRIns{
			Typ: codegen.CBR,
			Op:  codegen.BRNEQ,
			Arg1: &codegen.SymbolTableLiteralEntry{
				Value: 1,
				Repr:  "$1",
			},
			Arg2: v.expr.Symbol,
			Dst:  entry,
		})

		// If Body
		code = append(code, v.body.Code...)
		code = append(code, codegen.IRIns{
			Typ:  codegen.JMP,
			Op:   codegen.JMPO,
			Arg1: end,
		})
		code = append(code, codegen.IRIns{
			Typ: codegen.LBL,
			Dst: entry,
		})
	}

	code = append(code, elseBody.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: end,
	})

	ifElseCount++

	addrCode := &AddrCode{
		Code: code,
	}
	return addrCode, nil
}
