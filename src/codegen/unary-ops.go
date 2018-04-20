package codegen

import (
	"fmt"
	"log"
)

func genUOpCode(ins IRIns, regs [3]registerResult) {
	spill(regs[2].Spill)
	switch ins.Op {
	case NEG:
		var valueString string
		if regs[0].Register == "" {
			valueString = fmt.Sprintf("$%d", ins.Arg1.(*LiteralEntry).Value)
		} else {
			load(regs[0], ins.Arg1)
			valueString = string(regs[0].Register)
		}
		Code += fmt.Sprintf("movl %s, %s\n", valueString, regs[2].Register)
		Code += fmt.Sprintf("negl %s\n", regs[2].Register)
	case BNOT:
		load(regs[0], ins.Arg1)
		Code += fmt.Sprintf("movl %s, %s\n", regs[0].Register, regs[2].Register)
		Code += fmt.Sprintf("notl %s\n", regs[2].Register)
	case VAL:
		load(regs[0], ins.Arg1)
		Code += fmt.Sprintf("movl 0(%s), %s\n", regs[0].Register, regs[2].Register)
	case ADDR:
		Code += fmt.Sprintf("lea %s, %s\n", ins.Arg1.(*VariableEntry).MemoryLocation, regs[2].Register)
	default:
		log.Fatalf("Unknown op code for unary op: %s", ins.Op)
	}

	dst := ins.Dst.(*VariableEntry)
	updateVariable(dst, regs[2].Register)
}
