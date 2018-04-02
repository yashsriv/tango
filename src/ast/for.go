package ast

import (
	"fmt"
	"tango/src/codegen"
)

type forHeader struct {
	preStatement  *AddrCode
	Expr          *AddrCode
	postStatement *AddrCode
}

var forCount = 0

// EvalForHeader evaluates a ForHeader
func EvalForHeader(a, b, c Attrib) (forHeader, error) {

	preStatement, ok := a.(*AddrCode)
	if !ok {
		return forHeader{}, fmt.Errorf("unable to type cast %v to *AddrCode", a)
	}
	expr, ok := b.(*AddrCode)
	if !ok {
		return forHeader{}, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}
	postStatement, ok := c.(*AddrCode)
	if !ok {
		return forHeader{}, fmt.Errorf("unable to type cast %v to *AddrCode", c)
	}

	return forHeader{preStatement, expr, postStatement}, nil
}

// EvalForBody evaluates a ForBody
func EvalForBody(a, b Attrib) (*AddrCode, error) {

	header, ok := a.(forHeader)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to forHeader", a)
	}
	body, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}

	// Labels
	start, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_for_%d_start", forCount))
	if err != nil {
		return nil, err
	}
	end, err := codegen.InsertToSymbolTable(fmt.Sprintf("#_for_%d_end", forCount))
	if err != nil {
		return nil, err
	}
	forCount++

	// PreStatement
	code := header.preStatement.Code

	// Label for start, i.e. Expression Check
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: start,
	})
	// Check Expression
	if header.Expr != nil {
		code = append(code, header.Expr.Code...)
		code = append(code, codegen.IRIns{
			Typ: codegen.CBR,
			Op:  codegen.BRNEQ,
			Arg1: &codegen.SymbolTableLiteralEntry{
				Value: 1,
				Repr:  "$1",
			},
			Arg2: header.Expr.Symbol,
			Dst:  end,
		})
	}

	// Main body
	code = append(code, body.Code...)
	// Update Statement
	code = append(code, header.postStatement.Code...)
	// Jump to start
	code = append(code, codegen.IRIns{
		Typ:  codegen.JMP,
		Op:   codegen.JMPO,
		Arg1: start,
	})
	// Label for end - condition becomes false
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: end,
	})

	addrCode := &AddrCode{
		Code: code,
	}
	return addrCode, nil
}
