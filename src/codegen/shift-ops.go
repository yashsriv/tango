package codegen

import (
	"fmt"
	"log"
)

func genSOpCode(ins IRIns, regs [3]registerResult) {

	load(regs[0], ins.Arg1)
	// if arg1 register and dst register are same, it will be spilled in the previous statement
	if regs[0].Register != regs[2].Register {
		spill(regs[2].Spill)
	}

	// There is no issue like in genBOpCode because we have guaranteed that dst register is not
	// equal to arg1 register or arg2 register
	Code += fmt.Sprintf("movl %s, %s\n", regs[0].Register, regs[2].Register)

	var valueString string
	if regs[1].Register == "" {
		valueString = fmt.Sprintf("$%d", ins.Arg2.(*LiteralEntry).Value)
	} else {
		load(regs[1], ins.Arg2)
		valueString = "%cl"
	}

	switch ins.Op {
	case BSL:
		Code += fmt.Sprintf("shl %s, %s\n", valueString, regs[2].Register)
	case BSR:
		Code += fmt.Sprintf("shr %s, %s\n", valueString, regs[2].Register)
	default:
		log.Fatalf("Unknown op code for shift op: %s", ins.Op)
	}

	dst := ins.Dst.(*VariableEntry)
	updateVariable(dst, regs[2].Register)
}
