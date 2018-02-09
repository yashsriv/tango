package codegen

import "fmt"

var Code string

func spill(entries map[*SymbolTableRegisterEntry]bool) {
	for entry := range entries {
		Code += fmt.Sprintf("movl %s, %s\n", addrDesc[entry].regLocation, entry.Register)

		delete(regDesc[addrDesc[entry].regLocation], entry)

		addrDesc[entry] = address{
			regLocation: "",
			memLocation: entry.Register,
		}
	}
}

func load(regres registerResult, memloc SymbolTableEntry) {
	// Spill the register if needed
	reg := regres.Register
	spill(regres.Spill)

	// Load the value onto the register
	// can be a virtual register or a constant
	if _memloc, isRegister := memloc.(*SymbolTableRegisterEntry); isRegister {
		Code += fmt.Sprintf("movl %s, %s\n", _memloc.Register, reg)
		regDesc[reg][_memloc] = true
		addrDesc[_memloc] = address{
			regLocation: reg,
			memLocation: _memloc.Register,
		}
	} else {
		Code += fmt.Sprintf("movl %s, %s\n", memloc.(*SymbolTableLiteralEntry).SymbolTableString(), reg)
	}
}

func bopCode(ins IRIns, regs [3]MachineRegister) {
	Code += fmt.Sprintf("movl %s, %s", regs[1], regs[0])
	switch ins.Op {
	case ADD:
		Code += fmt.Sprintf("add %s, %s", regs[2], regs[0])
	case SUB:
		Code += fmt.Sprintf("sub %s, %s", regs[2], regs[0])
	case MUL:
		Code += fmt.Sprintf("imul %s, %s", regs[2], regs[0])
	case BSL:
		Code += fmt.Sprintf("shl %s, %s", regs[2], regs[0])
	case BSR:
		Code += fmt.Sprintf("shr %s, %s", regs[2], regs[0])
	}
	addrDesc[ins.Dst.(*SymbolTableRegisterEntry)] = address{
		regLocation: regs[0],
		memLocation: "",
	}
}

func op(ins IRIns, regs [3]registerResult) {
	switch ins.Typ {
	case JMP:
		Code += fmt.Sprintf("jmp %s\n", ins.Arg1.(*SymbolTableTargetEntry).Target)
	case KEY:
		switch ins.Op {
		case CALL:
			Code += fmt.Sprintf("call %s\n", ins.Arg1.(*SymbolTableTargetEntry).Target)
		case PARAM:
			if arg1, isLit := ins.Arg1.(*SymbolTableLiteralEntry); isLit {
				Code += fmt.Sprintf("push %s\n", arg1.SymbolTableString())
			} else {
				// TODO: Instead of loading to register. If not in a register,
				// then directly push from memory otherwise push from register
				load(regs[0], ins.Arg1)
				Code += fmt.Sprintf("push %s\n", regs[0].Register)
			}
		case RET:
			Code += "ret\n"
		case HALT:
			Code += "call exit\n"
			// TODO: Add _fmt* to .data segment
		case PRINTINT:
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
			Code += fmt.Sprintf("push $_fmtint\n")
			Code += fmt.Sprintf("call printf\n")
		case PRINTCHAR:
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
			Code += fmt.Sprintf("push $_fmtchar\n")
			Code += fmt.Sprintf("call printf\n")
		case PRINTSTR:
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
			Code += fmt.Sprintf("push $_fmtstr\n")
			Code += fmt.Sprintf("call printf\n")
		case SCANINT:
			Code += fmt.Sprintf("push $%s\n", ins.Arg1.(*SymbolTableRegisterEntry).Register)
			Code += fmt.Sprintf("push $_fmtint\n")
			Code += fmt.Sprintf("call scanf\n")
		case SCANCHAR:
			Code += fmt.Sprintf("push $%s\n", ins.Arg1.(*SymbolTableRegisterEntry).Register)
			Code += fmt.Sprintf("push $_fmtchar\n")
			Code += fmt.Sprintf("call scanf\n")
		case SCANSTR:
			Code += fmt.Sprintf("push $%s\n", ins.Arg1.(*SymbolTableRegisterEntry).Register)
			Code += fmt.Sprintf("push $_fmtstr\n")
			Code += fmt.Sprintf("call scanf\n")
		case INC:
			arg1 := ins.Arg1.(*SymbolTableRegisterEntry)
			if addrDesc[arg1].regLocation == "" {
				Code += fmt.Sprintf("inc $%s\n", arg1.Register)
			} else {
				Code += fmt.Sprintf("inc %s\n", addrDesc[arg1].regLocation)
				addrDesc[arg1] = address{
					regLocation: addrDesc[arg1].regLocation,
					memLocation: "",
				}
			}
		case DEC:
			arg1 := ins.Arg1.(*SymbolTableRegisterEntry)
			if addrDesc[arg1].regLocation == "" {
				Code += fmt.Sprintf("dec $%s\n", arg1.Register)
			} else {
				Code += fmt.Sprintf("dec %s\n", addrDesc[arg1].regLocation)
				addrDesc[arg1] = address{
					regLocation: addrDesc[arg1].regLocation,
					memLocation: "",
				}
			}
		}
	case ASN:
		load(regs[0], ins.Arg1)
		regDesc[regs[0].Register][ins.Dst.(*SymbolTableRegisterEntry)] = true
		addrDesc[ins.Dst.(*SymbolTableRegisterEntry)] = address{
			regLocation: regs[0].Register,
			memLocation: "",
		}
	case UOP:
		switch ins.Op {
		case NEG:
			load(regs[0], ins.Arg1)
			spill(regs[2].Spill)
			Code += fmt.Sprintf("movl %s, %s\n", regs[0].Register, regs[2].Register)
			Code += fmt.Sprintf("neg %s\n", regs[2].Register)
			regDesc[regs[2].Register][ins.Dst.(*SymbolTableRegisterEntry)] = true
			addrDesc[ins.Dst.(*SymbolTableRegisterEntry)] = address{
				regLocation: regs[2].Register,
				memLocation: "",
			}
		case BNOT:
			load(regs[0], ins.Arg1)
			spill(regs[2].Spill)
			Code += fmt.Sprintf("movl %s, %s\n", regs[0].Register, regs[2].Register)
			Code += fmt.Sprintf("not %s\n", regs[2].Register)
			regDesc[regs[2].Register][ins.Dst.(*SymbolTableRegisterEntry)] = true
			addrDesc[ins.Dst.(*SymbolTableRegisterEntry)] = address{
				regLocation: regs[2].Register,
				memLocation: "",
			}
		case VAL:
			// TODO: Discuss with Sir
			// Maintain a pointer map while dereferencing. Check if what we
			// want to dereference is in registers
			load(regs[0], ins.Arg1)
			spill(regs[2].Spill)
			Code += fmt.Sprintf("movl (%s), %s\n", regs[0].Register, regs[2].Register)
			regDesc[regs[2].Register][ins.Dst.(*SymbolTableRegisterEntry)] = true
			addrDesc[ins.Dst.(*SymbolTableRegisterEntry)] = address{
				regLocation: regs[2].Register,
				memLocation: "",
			}
		case ADDR:
			// TODO: Discuss with Sir
			// Maintain a pointer map while dereferencing. Check if what we
			// want to dereference is in registers
		}
	case CBR:
	}
}

func genCode() {
	for _, bbl := range BBLList {
		for i, ins := range bbl.Block {
			arg1res, arg2res, dstres := getReg(ins, bbl.Info[i])
			op(ins, [3]registerResult{arg1res, arg2res, dstres})

			// if arg1res.Register != "" && regDesc[arg1res.Register] != ins.Arg1 {
			// 	load(arg1res.Register, ins.Arg1)
			// }
			// if arg2, isRegister := ins.Arg2.(*SymbolTableRegisterEntry); isRegister && regDesc[arg2res.Register] != ins.Arg2 {
			// 	load(arg2res.Register, arg2)
			// }
			// if dst, isRegister := ins.Dst.(*SymbolTableRegisterEntry); isRegister && regDesc[dstres.Register] != ins.Dst {
			// 	load(dstres.Register, dst)
			// }

		}
	}
}
