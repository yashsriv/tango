package codegen

import (
	"fmt"
	"log"
)

func genBOpCode(ins IRIns, regs [3]registerResult) {

	load(regs[0], ins.Arg1)

	// if arg1 register and dst register are same, it will be spilled in the previous statement
	if regs[0].Register != regs[2].Register {
		spill(regs[2].Spill)
	}

	var valueString string
	if regs[1].Register == "" {
		valueString = fmt.Sprintf("$%d", ins.Arg2.(*LiteralEntry).Value)
	} else {
		load(regs[1], ins.Arg2)
		valueString = string(regs[1].Register)
	}

	var op1 string

	// dst reg == arg2 reg
	if regs[2].Register == regs[1].Register {
		// This has to be a SymbolTableVariableEntry since a register was allocated
		// if not program should fail... fatal error
		entry := ins.Arg2.(*VariableEntry)
		delete(regDesc[regs[1].Register], entry)
		addrDesc[entry] = address{
			regLocation: "",
			memLocation: entry.MemoryLocation,
		}
		op1 = string(regs[0].Register)
	} else if regs[2].Register == regs[0].Register {
		// dst reg == arg1 reg
		// This has to be a SymbolTableVariableEntry since a register was allocated
		// if not program should fail... fatal error
		entry := ins.Arg1.(*VariableEntry)
		delete(regDesc[regs[0].Register], entry)
		addrDesc[entry] = address{
			regLocation: "",
			memLocation: entry.MemoryLocation,
		}
		op1 = valueString
	} else {
		// All registers unique
		// Move arg1 register to dst register
		Code += fmt.Sprintf("movl %s, %s\n", regs[0].Register, regs[2].Register)
		op1 = valueString
	}

	switch ins.Op {
	case ADD:
		Code += fmt.Sprintf("addl %s, %s\n", op1, regs[2].Register)
	case SUB:
		Code += fmt.Sprintf("subl %s, %s\n", op1, regs[2].Register)
	case MUL:
		Code += fmt.Sprintf("imul %s, %s\n", op1, regs[2].Register)
	case AND:
		fallthrough
	case BAND:
		Code += fmt.Sprintf("andl %s, %s\n", op1, regs[2].Register)
	case OR:
		fallthrough
	case BOR:
		Code += fmt.Sprintf("orl %s, %s\n", op1, regs[2].Register)
	case XOR:
		Code += fmt.Sprintf("xorl %s, %s\n", op1, regs[2].Register)
	default:
		log.Fatalf("Unknown op code for binary op: %s", ins.Op)
	}

	dst := ins.Dst.(*VariableEntry)
	updateVariable(dst, regs[2].Register)
}
