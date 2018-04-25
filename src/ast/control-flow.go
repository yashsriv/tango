package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

type labelName struct {
	start *codegen.TargetEntry
	end   *codegen.TargetEntry
}

// EvalLabelName is used to evaluate a label
func EvalLabelName(a Attrib) (labelName, error) {
	identifier := string(a.(*token.Token).Lit)
	start := &codegen.TargetEntry{
		Target: fmt.Sprintf("_label_%s", identifier),
	}
	// Associating identifier with some entry
	err := codegen.SymbolTable.InsertSymbol(identifier, start)
	end := &codegen.TargetEntry{
		Target: fmt.Sprintf("_labelend_%s", identifier),
	}
	return labelName{start, end}, err
}

// EvalLabel is used to evaluate a label
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

var currentRetType codegen.TypeEntry

// EvalReturn evaluates a return statement
func EvalReturn(a Attrib) (*AddrCode, error) {
	expr, ok := a.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[EvalReturn] unable to typecast %v to []*AddrCode", a)
	}
	code := make([]codegen.IRIns, 0)
	switch len(expr) {
	case 0:
		if currentRetType != VoidType {
			return nil, fmt.Errorf("Expected a return value of type: %v", currentRetType)
		}
		code = append(code, codegen.IRIns{
			Typ: codegen.KEY,
			Op:  codegen.RET,
		})
	case 1:
		if !SameType(currentRetType, expr[0].Symbol.Type()) {
			return nil, fmt.Errorf("Expected a return value of type: %v", currentRetType)
		}
		evalExpr := EvalWrapped(expr[0])
		code = append(code, evalExpr.Code...)
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.RETI,
			Arg1: evalExpr.Symbol,
		})
	default:
		return nil, ErrUnsupported
	}
	return &AddrCode{
		Code: code,
	}, nil
}

// EvalCall evaluates a call statement
func EvalCall(a, b Attrib) (*AddrCode, error) {
	entry_ := a.(*AddrCode).Symbol
	entry, ok := entry_.(*codegen.TargetEntry)
	var varEntry *codegen.VariableEntry
	if !ok {
		if variable, ok1 := entry_.(*codegen.VariableEntry); ok1 {
			if t, ok2 := variable.Type().(codegen.FuncType); ok2 {
				entry = t.Target
				varEntry = variable
			} else {
				return nil, fmt.Errorf("invalid function call statement")
			}
		} else {
			return nil, fmt.Errorf("invalid function call statement")
		}
	}
	exprList, ok := b.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to []*AddrCode", b)
	}

	var isSystem bool

	if entry.Target == "printf" || entry.Target == "scanf" {
		isSystem = true
	}

	if !isSystem && len(exprList) != len(entry.InType) {
		return nil, fmt.Errorf("wrong number of arguments in function call. expected %d, got %d", len(entry.InType), len(exprList))
	}
	code := a.(*AddrCode).Code
	for i := len(exprList) - 1; i >= 0; i-- {
		if !isSystem && !SameType(exprList[i].Symbol.Type(), entry.InType[i]) {
			return nil, fmt.Errorf("wrong type of argument %d in function call. expected %v, got %v", i, entry.InType[i], exprList[i].Symbol.Type())
		}
		if isSystem && i == 0 {
			if !SameType(exprList[i].Symbol.Type(), stringType) {
				return nil, fmt.Errorf("wrong type of 1st argument in function call. expected %v, got %v", stringType, exprList[i].Symbol.Type())
			}
		}
		evaluatedExpr := EvalWrapped(exprList[i])
		code = append(code, evaluatedExpr.Code...)
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.PARAM,
			Arg1: evaluatedExpr.Symbol,
		})
	}
	if varEntry == nil {
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.CALL,
			Arg1: entry,
		})
	} else {
		code = append(code, codegen.IRIns{
			Typ:  codegen.KEY,
			Op:   codegen.CALL,
			Arg1: varEntry,
		})
	}
	// code = append(code, codegen.IRIns{
	// 	Typ: codegen.KEY,
	// 	Op:  codegen.UNALLOC,
	// 	Arg1: &codegen.LiteralEntry{
	// 		Value: len(exprList) * 4,
	// 		LType: intType,
	// 	},
	// })
	entry1 := CreateTemporary(entry.RetType)
	code = append(code, entry1.Code...)
	code = append(code, codegen.IRIns{
		Typ:  codegen.KEY,
		Op:   codegen.SETRET,
		Arg1: entry1.Symbol,
	})
	return &AddrCode{
		Code:   code,
		Symbol: entry1.Symbol,
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
