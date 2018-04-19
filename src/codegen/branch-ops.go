package codegen

import (
	"fmt"
	"log"
)

func genCBRCode(ins IRIns, regs [3]registerResult) {
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
	case BREQ:
		Code += fmt.Sprintf("je %s\n", ins.Dst.(*TargetEntry).Target)
	case BRNEQ:
		Code += fmt.Sprintf("jne %s\n", ins.Dst.(*TargetEntry).Target)
	case BRLT:
		Code += fmt.Sprintf("jl %s\n", ins.Dst.(*TargetEntry).Target)
	case BRLTE:
		Code += fmt.Sprintf("jle %s\n", ins.Dst.(*TargetEntry).Target)
	case BRGT:
		Code += fmt.Sprintf("jg %s\n", ins.Dst.(*TargetEntry).Target)
	case BRGTE:
		Code += fmt.Sprintf("jge %s\n", ins.Dst.(*TargetEntry).Target)
	default:
		log.Fatalf("Unknown op code for branch op: %s", ins.Op)
	}
}
