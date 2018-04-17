package ast

import (
	"fmt"
	"tango/src/codegen"
)

type forHeader struct {
	preStatement  *AddrCode
	Expr          *AddrCode
	postStatement *AddrCode
	start         *codegen.SymbolTableTargetEntry
	end           *codegen.SymbolTableTargetEntry
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

	start := &codegen.SymbolTableTargetEntry{
		Target: fmt.Sprintf("_for_%d_start", forCount),
	}
	end := &codegen.SymbolTableTargetEntry{
		Target: fmt.Sprintf("_for_%d_end", forCount),
	}
	continueStack = continueStack.Push(start)
	breakStack = breakStack.Push(end)
	forCount++
	return forHeader{preStatement, expr, postStatement, start, end}, nil
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
	start := header.start
	end := header.end
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

	var entry *codegen.SymbolTableTargetEntry
	breakStack, entry = breakStack.Pop()
	if entry != end {
		return nil, fmt.Errorf("break labels do not match. something is very very wrong")
	}
	continueStack, entry = continueStack.Pop()
	if entry != start {
		return nil, fmt.Errorf("continue labels do not match. something is very very wrong")
	}
	return addrCode, nil
}
