package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

type labelName struct {
	start *codegen.SymbolTableTargetEntry
	end   *codegen.SymbolTableTargetEntry
}

func EvalLabelName(a Attrib) (labelName, error) {
	identifier := string(a.(*token.Token).Lit)
	if _, ok := codegen.AccSymbolMap(identifier); ok {
		return labelName{}, fmt.Errorf("Identifier %s is being declared twice in this scope", identifier)
	}
	start := &codegen.SymbolTableTargetEntry{
		Target: fmt.Sprintf("_label_%s", identifier),
	}
	codegen.InsertToSymbolMap(identifier, start)
	end := &codegen.SymbolTableTargetEntry{
		Target: fmt.Sprintf("_labelend_%s", identifier),
	}
	return labelName{start, end}, nil
}

func EvalLabel(a, b Attrib) (*AddrCode, error) {
	label, ok := a.(labelName)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to labelName", a)
	}
	stmt, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", b)
	}

	code := make([]codegen.IRIns, 0)
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: label.start,
	})
	code = append(code, stmt.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: label.end,
	})
	return &AddrCode{
		Code: code,
	}, nil
}

func EvalGoto(a Attrib) (*AddrCode, error) {
	identifier := string(a.(*token.Token).Lit)
	entry, ok := codegen.AccSymbolMap(identifier)
	if !ok {
		return nil, fmt.Errorf("Identifier %s is undefined in this scope", identifier)
	}
	code := make([]codegen.IRIns, 0)
	code = append(code, codegen.IRIns{
		Typ:  codegen.JMP,
		Op:   codegen.JMPO,
		Arg1: entry,
	})
	return &AddrCode{
		Code: code,
	}, nil
}

func EvalReturn(a Attrib) (*AddrCode, error) {
	expr, ok := a.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", a)
	}
	code := make([]codegen.IRIns, 0)
	switch len(expr) {
	case 0:
		code = append(code, codegen.IRIns{
			Typ: codegen.KEY,
			Op:  codegen.RET,
		})
	case 1:
		code = append(code, expr[0].Code...)
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.RETI,
			Arg1: expr[0].Symbol,
		})
	default:
		return nil, ErrUnsupported
	}
	return &AddrCode{
		Code: code,
	}, nil
}

func EvalCall(a, b Attrib) (*AddrCode, error) {
	entry_ := a.(*AddrCode).Symbol
	entry, ok := entry_.(*codegen.SymbolTableTargetEntry)
	if !ok {
		return nil, fmt.Errorf("invalid function call statement")
	}
	exprList, ok := b.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to []*AddrCode", b)
	}
	code := make([]codegen.IRIns, 0)
	for i := len(exprList) - 1; i >= 0; i-- {
		code = append(code, exprList[i].Code...)
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.ARG,
			Arg1: exprList[i].Symbol,
		})
	}
	code = append(code, codegen.IRIns{
		Typ:  codegen.KEY,
		Op:   codegen.CALL,
		Arg1: entry,
	})
	entry1, err := codegen.InsertToSymbolTable(fmt.Sprintf("rtmp%d", tempCount))
	if err != nil {
		return nil, err
	}
	code = append(code, codegen.IRIns{
		Typ:  codegen.KEY,
		Op:   codegen.SETRET,
		Arg1: entry1,
	})
	tempCount++
	return &AddrCode{
		Code:   code,
		Symbol: entry1,
	}, nil
}
