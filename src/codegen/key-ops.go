package codegen

import (
	"fmt"
	"log"
)

func genKeyCode(ins IRIns, regs [3]registerResult) {
	switch ins.Op {
	case CALL:
		Code += fmt.Sprintf("call %s\n", ins.Arg1.(*TargetEntry).Target)
	case PARAM:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push $%d\n", ins.Arg1.(*LiteralEntry).Value)
		} else {
			// TODO: Instead of loading to register. If not in a register,
			// then directly push from memory otherwise push from register
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
		}
	case SETRET:
		regDesc[returnRegister][ins.Arg1.(*VariableEntry)] = true
		addrDesc[ins.Arg1.(*VariableEntry)] = address{
			regLocation: returnRegister,
			memLocation: nil,
		}
	case RETI:
		load(registerResult{Register: returnRegister}, ins.Arg1)
		fallthrough
	case RET:
		Code += "movl %ebp, %esp\n"
		Code += "pop %ebp\n"
		Code += "ret\n"
	case HALT:
		Code += "call exit\n"
	case PRINTINT:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push $%d\n", ins.Arg1.(*LiteralEntry).Value)
		} else {
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
		}
		Code += fmt.Sprintf("push $_fmtint\n")
		Code += fmt.Sprintf("call printf\n")
	case PRINTCHAR:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push $%d\n", ins.Arg1.(*LiteralEntry).Value)
		} else {
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
		}
		Code += fmt.Sprintf("push $_fmtchar\n")
		Code += fmt.Sprintf("call printf\n")
	case PRINTSTR:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push $%d\n", ins.Arg1.(*LiteralEntry).Value)
		} else {
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
		}
		Code += fmt.Sprintf("push $_fmtstr\n")
		Code += fmt.Sprintf("call printf\n")
	case SCANINT:

		arg1 := ins.Arg1.(*VariableEntry)
		Code += fmt.Sprintf("lea %s, %%eax\n", arg1.MemoryLocation)
		Code += fmt.Sprintf("push %%eax\n")
		Code += fmt.Sprintf("push $_fmtint\n")
		Code += fmt.Sprintf("call scanf\n")

		// Invalidate registers if any
		if val, ok := addrDesc[arg1]; ok && val.regLocation != "" {
			delete(regDesc[val.regLocation], arg1)
			addrDesc[arg1] = address{
				regLocation: "",
				memLocation: arg1.MemoryLocation,
			}
		} else if !ok {
			log.Fatalf("[scanint] addrDesc missing")
		}

	case SCANCHAR:
		arg1 := ins.Arg1.(*VariableEntry)
		Code += fmt.Sprintf("lea %s, %%eax\n", arg1.MemoryLocation)
		Code += fmt.Sprintf("push %%eax\n")
		Code += fmt.Sprintf("push $_fmtchar\n")
		Code += fmt.Sprintf("call scanf\n")
		// Invalidate registers if any
		if val, ok := addrDesc[arg1]; ok && val.regLocation != "" {
			delete(regDesc[val.regLocation], arg1)
			addrDesc[arg1] = address{
				regLocation: "",
				memLocation: arg1.MemoryLocation,
			}
		} else if !ok {
			log.Fatalf("[scanint] addrDesc missing")
		}

	case SCANSTR:
		arg1 := ins.Arg1.(*VariableEntry)
		Code += fmt.Sprintf("lea %s, %%eax\n", arg1.MemoryLocation)
		Code += fmt.Sprintf("push %%eax\n")
		Code += fmt.Sprintf("push $_fmtstr\n")
		Code += fmt.Sprintf("call scanf\n")
		// Invalidate registers if any
		if val, ok := addrDesc[arg1]; ok && val.regLocation != "" {
			delete(regDesc[val.regLocation], arg1)
			addrDesc[arg1] = address{
				regLocation: "",
				memLocation: arg1.MemoryLocation,
			}
		} else if !ok {
			log.Fatalf("[scanint] addrDesc missing")
		}

	case INC:
		arg1 := ins.Arg1.(*VariableEntry)
		if val, ok := addrDesc[arg1]; ok && val.regLocation == "" {
			Code += fmt.Sprintf("incl %s\n", arg1.MemoryLocation)
		} else if ok {
			Code += fmt.Sprintf("incl %s\n", val.regLocation)
			addrDesc[arg1] = address{
				regLocation: val.regLocation,
				memLocation: nil,
			}
		} else {
			log.Fatalf("[INC] AddrDesc missing")
		}
	case DEC:
		arg1 := ins.Arg1.(*VariableEntry)
		if val, ok := addrDesc[arg1]; ok && val.regLocation == "" {
			Code += fmt.Sprintf("decl %s\n", arg1.MemoryLocation)
		} else if ok {
			Code += fmt.Sprintf("decl %s\n", val.regLocation)
			addrDesc[arg1] = address{
				regLocation: val.regLocation,
				memLocation: nil,
			}
		} else {
			log.Fatalf("[DEC] AddrDesc missing")
		}
	case ALLOC:
		arg1 := ins.Arg1.(*LiteralEntry)
		Code += fmt.Sprintf("sub %s, %%esp\n", arg1)
	case UNALLOC:
		arg1 := ins.Arg1.(*LiteralEntry)
		Code += fmt.Sprintf("add %s, %%esp\n", arg1)
	case TAKE:
		// Dst where to write the value
		load(regs[2], ins.Dst)
		// From which memory location we want to read the value
		load(regs[0], ins.Arg1)

		if regs[0].Register == "" || regs[2].Register == "" {
			log.Fatalf("[TAKE] cannot have empty register")
		}

		var op1 string
		// Memory offset
		if regs[1].Register != "" {
			load(regs[1], ins.Arg2)
			op1 = string(regs[1].Register)
		}
		if op1 == "" {
			offset := ins.Arg2.(*LiteralEntry).Value * 4
			Code += fmt.Sprintf("movl %d(%s), %s\n", offset, regs[0].Register, regs[2].Register)
		} else {
			Code += fmt.Sprintf("movl (%s, %s, 4), %s\n", regs[0].Register, op1, regs[2].Register)
		}
		dst := ins.Dst.(*VariableEntry)
		updateVariable(dst, regs[2].Register)
	case PUT:
		load(regs[2], ins.Dst)
		var op1, op2 string
		if regs[0].Register != "" {
			load(regs[0], ins.Arg1)
			op1 = string(regs[0].Register)
		}
		if regs[1].Register == "" {
			op2 = fmt.Sprintf("$%d", ins.Arg2.(*LiteralEntry).Value)
		} else {
			load(regs[1], ins.Arg2)
			op2 = string(regs[1].Register)
		}
		if op1 == "" {
			offset := ins.Arg1.(*LiteralEntry).Value * 4
			Code += fmt.Sprintf("movl %s, %d(%s)\n", op2, offset, regs[2].Register)
		} else {
			Code += fmt.Sprintf("movl %s, (%s, %s, 4)\n", op2, regs[2].Register, op1)
		}
		dst := ins.Dst.(*VariableEntry)
		updateVariable(dst, regs[2].Register)
	case MALLOC:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push $%d\n", ins.Arg1.(*LiteralEntry).Value*4)
		} else {
			log.Fatalf("non constant value not supported for malloc ops")
		}
		Code += fmt.Sprintf("call malloc\n")

	default:
		log.Fatalf("Unknown op code for key op: %s", ins.Op)
	}
}
