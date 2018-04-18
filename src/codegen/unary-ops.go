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
		log.Fatalf("Unhandled pointer stuff")
		// TODO: Discuss with Sir
		// Maintain a pointer map while dereferencing. Check if what we
		// want to dereference is in registers
	case ADDR:
		log.Fatalf("Unhandled pointer stuff")
		// TODO: Discuss with Sir
		// Maintain a pointer map while dereferencing. Check if what we
		// want to dereference is in registers
	}

	dst := ins.Dst.(*VariableEntry)
	updateVariable(dst, regs[2].Register)
}
