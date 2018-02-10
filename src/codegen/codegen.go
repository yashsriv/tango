package codegen

import (
	"fmt"
	"log"
)

var Code string

func spill(entries map[*SymbolTableVariableEntry]bool) {
	for entry := range entries {
		Code += fmt.Sprintf("movl %s, (%s)\n", addrDesc[entry].regLocation, entry.MemoryLocation)

		delete(regDesc[addrDesc[entry].regLocation], entry)

		addrDesc[entry] = address{
			regLocation: "",
			memLocation: entry.MemoryLocation,
		}
	}
}

func load(regres registerResult, memloc SymbolTableEntry) {

	// Spill the register if needed
	reg := regres.Register

	if reg == "" {
		log.Fatalf("Trying to load into an empty register\n")
	}

	if _memloc, isRegister := memloc.(*SymbolTableVariableEntry); isRegister {
		// If we are trying to shift a variable already in the register, ignore
		if _, isKey := regDesc[reg][_memloc]; isKey {
			return
		}
	}

	spill(regres.Spill)

	// Load the value onto the register
	// can be a virtual register or a constant
	if _memloc, isRegister := memloc.(*SymbolTableVariableEntry); isRegister {
		Code += fmt.Sprintf("movl (%s), %s\n", _memloc.MemoryLocation, reg)
		regDesc[reg][_memloc] = true
		addrDesc[_memloc] = address{
			regLocation: reg,
			memLocation: _memloc.MemoryLocation,
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
	addrDesc[ins.Dst.(*SymbolTableVariableEntry)] = address{
		regLocation: regs[0],
		memLocation: "",
	}
}

func genKeyCode(ins IRIns, regs [3]registerResult) {
	switch ins.Op {
	case CALL:
		Code += fmt.Sprintf("call %s\n", ins.Arg1.(*SymbolTableTargetEntry).Target)
	case PARAM:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push %s\n", ins.Arg1.SymbolTableString())
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
	case PRINTINT:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push %s\n", ins.Arg1.SymbolTableString())
		} else {
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
		}
		Code += fmt.Sprintf("push $_fmtint\n")
		Code += fmt.Sprintf("call printf\n")
	case PRINTCHAR:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push %s\n", ins.Arg1.SymbolTableString())
		} else {
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
		}
		Code += fmt.Sprintf("push $_fmtchar\n")
		Code += fmt.Sprintf("call printf\n")
	case PRINTSTR:
		if regs[0].Register == "" {
			// This is a literal. Push directly.
			Code += fmt.Sprintf("push %s\n", ins.Arg1.SymbolTableString())
		} else {
			load(regs[0], ins.Arg1)
			Code += fmt.Sprintf("push %s\n", regs[0].Register)
		}
		Code += fmt.Sprintf("push $_fmtstr\n")
		Code += fmt.Sprintf("call printf\n")
	case SCANINT:

		arg1 := ins.Arg1.(*SymbolTableVariableEntry)
		Code += fmt.Sprintf("push $%s\n", arg1.MemoryLocation)
		Code += fmt.Sprintf("push $_fmtint\n")
		Code += fmt.Sprintf("call scanf\n")

		// Invalidate registers if any
		if addrDesc[arg1].regLocation != "" {
			delete(regDesc[addrDesc[arg1].regLocation], arg1)
			addrDesc[arg1] = address{
				regLocation: "",
				memLocation: arg1.MemoryLocation,
			}
		}

	case SCANCHAR:
		arg1 := ins.Arg1.(*SymbolTableVariableEntry)
		Code += fmt.Sprintf("push $%s\n", arg1.MemoryLocation)
		Code += fmt.Sprintf("push $_fmtchar\n")
		Code += fmt.Sprintf("call scanf\n")
		// Invalidate registers if any
		if addrDesc[arg1].regLocation != "" {
			delete(regDesc[addrDesc[arg1].regLocation], arg1)
			addrDesc[arg1] = address{
				regLocation: "",
				memLocation: arg1.MemoryLocation,
			}
		}

	case SCANSTR:
		arg1 := ins.Arg1.(*SymbolTableVariableEntry)
		Code += fmt.Sprintf("push $%s\n", arg1.MemoryLocation)
		Code += fmt.Sprintf("push $_fmtstr\n")
		Code += fmt.Sprintf("call scanf\n")
		// Invalidate registers if any
		if addrDesc[arg1].regLocation != "" {
			delete(regDesc[addrDesc[arg1].regLocation], arg1)
			addrDesc[arg1] = address{
				regLocation: "",
				memLocation: arg1.MemoryLocation,
			}
		}

	case INC:
		arg1 := ins.Arg1.(*SymbolTableVariableEntry)
		if addrDesc[arg1].regLocation == "" {
			Code += fmt.Sprintf("incl (%s)\n", arg1.MemoryLocation)
		} else {
			Code += fmt.Sprintf("incl %s\n", addrDesc[arg1].regLocation)
			addrDesc[arg1] = address{
				regLocation: addrDesc[arg1].regLocation,
				memLocation: "",
			}
		}
	case DEC:
		arg1 := ins.Arg1.(*SymbolTableVariableEntry)
		if addrDesc[arg1].regLocation == "" {
			Code += fmt.Sprintf("decl (%s)\n", arg1.MemoryLocation)
		} else {
			Code += fmt.Sprintf("decl %s\n", addrDesc[arg1].regLocation)
			addrDesc[arg1] = address{
				regLocation: addrDesc[arg1].regLocation,
				memLocation: "",
			}
		}
	}
}

func genUOpCode(ins IRIns, regs [3]registerResult) {
	switch ins.Op {
	case NEG:
		var valueString string
		if regs[0].Register == "" {
			valueString = ins.Arg1.SymbolTableString()
		} else {
			load(regs[0], ins.Arg1)
			valueString = string(regs[0].Register)
		}
		spill(regs[2].Spill)
		Code += fmt.Sprintf("movl %s, %s\n", valueString, regs[2].Register)
		Code += fmt.Sprintf("negl %s\n", regs[2].Register)
		regDesc[regs[2].Register][ins.Dst.(*SymbolTableVariableEntry)] = true
		addrDesc[ins.Dst.(*SymbolTableVariableEntry)] = address{
			regLocation: regs[2].Register,
			memLocation: "",
		}
	case BNOT:
		load(regs[0], ins.Arg1)
		spill(regs[2].Spill)
		Code += fmt.Sprintf("movl %s, %s\n", regs[0].Register, regs[2].Register)
		Code += fmt.Sprintf("notl %s\n", regs[2].Register)
		regDesc[regs[2].Register][ins.Dst.(*SymbolTableVariableEntry)] = true
		addrDesc[ins.Dst.(*SymbolTableVariableEntry)] = address{
			regLocation: regs[2].Register,
			memLocation: "",
		}
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
}

func genCBRCode(ins IRIns, regs [3]registerResult) {
	var op1, op2 string
	if regs[0].Register == "" {
		op1 = ins.Arg1.SymbolTableString()
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
		Code += fmt.Sprintf("je %s\n", ins.Dst.(*SymbolTableTargetEntry).Target)
	case BRNEQ:
		Code += fmt.Sprintf("jne %s\n", ins.Dst.(*SymbolTableTargetEntry).Target)
	case BRLT:
		Code += fmt.Sprintf("jl %s\n", ins.Dst.(*SymbolTableTargetEntry).Target)
	case BRLTE:
		Code += fmt.Sprintf("jle %s\n", ins.Dst.(*SymbolTableTargetEntry).Target)
	case BRGT:
		Code += fmt.Sprintf("jg %s\n", ins.Dst.(*SymbolTableTargetEntry).Target)
	case BRGTE:
		Code += fmt.Sprintf("jge %s\n", ins.Dst.(*SymbolTableTargetEntry).Target)
	}
}

func genOpCode(ins IRIns, pointerMap map[*SymbolTableVariableEntry]*SymbolTableVariableEntry, regs [3]registerResult) {
	switch ins.Typ {
	case LBL:
		Code += fmt.Sprintf("%s:\n", ins.Dst.(*SymbolTableTargetEntry).Target)
	case JMP:
		Code += fmt.Sprintf("jmp %s\n", ins.Arg1.(*SymbolTableTargetEntry).Target)
	case KEY:
		genKeyCode(ins, regs)
	case ASN:
		load(regs[0], ins.Arg1)
		regDesc[regs[0].Register][ins.Dst.(*SymbolTableVariableEntry)] = true
		addrDesc[ins.Dst.(*SymbolTableVariableEntry)] = address{
			regLocation: regs[0].Register,
			memLocation: "",
		}
	case UOP:
		genUOpCode(ins, regs)
	case CBR:
		genCBRCode(ins, regs)
	default:
		log.Fatalf("Unhandled instruction: %s\n", ins.String())
	}
}

func genCode() {
	pointerMap := make(map[*SymbolTableVariableEntry]*SymbolTableVariableEntry)

	for _, bbl := range BBLList {
		Code += "# Begin Basic Block\n"
		for i, ins := range bbl.Block {
			arg1res, arg2res, dstres := getReg(ins, bbl.Info[i])

			if i == len(bbl.Block)-1 && shouldPreSave(ins) {
				saveBBL()
			}

			genOpCode(ins, pointerMap, [3]registerResult{arg1res, arg2res, dstres})

			if i == len(bbl.Block)-1 && !shouldPreSave(ins) {
				saveBBL()
			}

		}
		for register := range regDesc {
			regDesc[register] = make(map[*SymbolTableVariableEntry]bool)
		}
		Code += "# End Basic Block\n"
	}
}

func genData() {
	Code += ".section .data\n\n"
	Code += "_fmtint: .string \"%d\"\n"
	Code += "_fmtchar: .string \"%c\"\n"
	Code += "_fmtstr: .string \"%s\"\n"
	for _, symbol := range SymbolTable {
		if variable, isVar := symbol.(*SymbolTableVariableEntry); isVar {
			Code += fmt.Sprintf("%s: .long 0\n", variable.MemoryLocation)
		}
	}
}

func saveBBL() {
	// Code += "\n\n# Saving Stuff\n"
	for register, variables := range regDesc {
		for variable := range variables {
			if addrDesc[variable].memLocation == "" {
				Code += fmt.Sprintf("movl %s, (%s)\n", register, variable.MemoryLocation)
			}
			addrDesc[variable] = address{
				regLocation: "",
				memLocation: variable.MemoryLocation,
			}
		}
	}
	// Code += "# Done Saving Stuff\n\n"
}

func genMisc() {
	Code += `.section .text
.globl main

`
}

// GenerateASM generates the assemby code in Code
func GenerateASM() {
	genMisc()
	genCode()
	genData()
}

func shouldPreSave(ins IRIns) bool {
	return ins.Typ == CBR || ins.Typ == JMP ||
		(ins.Typ == KEY && (ins.Op != PARAM && ins.Op != INC && ins.Op != DEC))
}
