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

// EvalLabelName is used to evaluate a label
func EvalLabelName(a Attrib) (labelName, error) {
	identifier := string(a.(*token.Token).Lit)
	start := &codegen.SymbolTableTargetEntry{
		Target: fmt.Sprintf("_label_%s", identifier),
	}
	// Associating identifier with some entry
	err := codegen.SymbolTable.InsertSymbol(identifier, start)
	end := &codegen.SymbolTableTargetEntry{
		Target: fmt.Sprintf("_labelend_%s", identifier),
	}
	return labelName{start, end}, err
}

func EvalLabel(a, b Attrib) (*AddrCode, error) {
	label, ok := a.(labelName)
	if !ok {
		return nil, fmt.Errorf("[EvalLabel] unable to type cast %v to labelName", a)
	}
	stmt, ok := b.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[EvalLabel] unable to type cast %v to *AddrCode", b)
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

// EvalGoto evaluates a goto statement
func EvalGoto(a Attrib) (*AddrCode, error) {
	identifier := string(a.(*token.Token).Lit)
	entry, err := codegen.SymbolTable.GetSymbol(identifier)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("[EvalReturn] unable to typecast %v to []*AddrCode", a)
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
	entry1 := CreateTemporary()
	code = append(code, codegen.IRIns{
		Typ:  codegen.KEY,
		Op:   codegen.SETRET,
		Arg1: entry1,
	})
	return &AddrCode{
		Code:   code,
		Symbol: entry1,
	}, nil
}

// NewScope marks the start of a scope
func NewScope() (Attrib, error) {
	codegen.NewScope()
	return nil, nil
}

// EndScope marks end of the scope
func EndScope() (Attrib, error) {
	err := codegen.EndScope()
	return nil, err
}
