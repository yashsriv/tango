package ast

import (
	"fmt"
	"tango/src/codegen"
)

type labelStack []*codegen.SymbolTableTargetEntry

func (s labelStack) Push(v *codegen.SymbolTableTargetEntry) labelStack {
	return append(s, v)
}

func (s labelStack) Pop() (labelStack, *codegen.SymbolTableTargetEntry) {
	// FIXME: What do we do if the labelStack is empty, though?

	l := len(s)
	return s[:l-1], s[l-1]
}

func (s labelStack) Empty() bool {
	return len(s) == 0
}

var breakStack labelStack
var continueStack labelStack

func EvalBreak() (*AddrCode, error) {
	if breakStack.Empty() {
		return nil, fmt.Errorf("misplaced break statement")
	}
	_, top := breakStack.Pop()

	code := make([]codegen.IRIns, 1)
	code[0] = codegen.IRIns{
		Typ:  codegen.JMP,
		Op:   codegen.JMPO,
		Arg1: top,
	}
	return &AddrCode{
		Symbol: nil,
		Code:   code,
	}, nil
}

func EvalContinue() (*AddrCode, error) {
	if continueStack.Empty() {
		return nil, fmt.Errorf("misplaced continue statement")
	}
	_, top := continueStack.Pop()

	code := make([]codegen.IRIns, 1)
	code[0] = codegen.IRIns{
		Typ:  codegen.JMP,
		Op:   codegen.JMPO,
		Arg1: top,
	}
	return &AddrCode{
		Symbol: nil,
		Code:   code,
	}, nil
}
