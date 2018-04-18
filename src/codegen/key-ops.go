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
		Code += fmt.Sprintf("push $%s\n", arg1.MemoryLocation)
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
		Code += fmt.Sprintf("push $%s\n", arg1.MemoryLocation)
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
		Code += fmt.Sprintf("push $%s\n", arg1.MemoryLocation)
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
			Code += fmt.Sprintf("incl (%s)\n", arg1.MemoryLocation)
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
			Code += fmt.Sprintf("decl (%s)\n", arg1.MemoryLocation)
		} else if ok {
			Code += fmt.Sprintf("decl %s\n", val.regLocation)
			addrDesc[arg1] = address{
				regLocation: val.regLocation,
				memLocation: nil,
			}
		} else {
			log.Fatalf("[DEC] AddrDesc missing")
		}
	}
}
