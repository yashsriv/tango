package codegen

import (
	"fmt"
	"log"
)

func genLOpCode(ins IRIns, regs [3]registerResult) {
	load(regs[0], ins.Arg1)
	spill(regs[2].Spill)

	var op1, op2 string
	if regs[0].Register == "" {
		op1 = fmt.Sprintf("$%d", ins.Arg1.(*LiteralEntry).Value)
	} else {
		load(regs[0], ins.Arg1)
		op1 = string(regs[0].Register)
	}
	if regs[1].Register == "" {
		log.Fatalf("2nd Operand of cmp cannot be a constant")
	} else {
		load(regs[1], ins.Arg2)
		op2 = string(regs[1].Register)
	}
	Code += fmt.Sprintf("cmpl %s, %s\n", op1, op2)

	switch ins.Op {
	case EQ:
		Code += fmt.Sprintf("je _logic_start_%d\n", logicCounter)
	case NEQ:
		Code += fmt.Sprintf("jne _logic_start_%d\n", logicCounter)
	case LT:
		Code += fmt.Sprintf("jl _logic_start_%d\n", logicCounter)
	case LTE:
		Code += fmt.Sprintf("jle _logic_start_%d", logicCounter)
	case GT:
		Code += fmt.Sprintf("jg _logic_start_%d", logicCounter)
	case GTE:
		Code += fmt.Sprintf("jge _logic_start_%d", logicCounter)
	default:
		log.Fatalf("Unknown op code for logical op: %s", ins.Op)
	}

	Code += fmt.Sprintf("movl $0, %s\n", regs[2].Register)
	Code += fmt.Sprintf("jmp _logic_end_%d", logicCounter)
	Code += fmt.Sprintf("_logic_start_%d:", logicCounter)
	Code += fmt.Sprintf("movl $1, %s\n", regs[2].Register)
	Code += fmt.Sprintf("_logic_end_%d:", logicCounter)
	logicCounter++

	dst := ins.Dst.(*VariableEntry)
	updateVariable(dst, regs[2].Register)
}
