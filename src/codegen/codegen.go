package codegen

import (
	"fmt"
	"log"
	"strings"
)

// Code stores the assembly of our IR.
var Code string
var logicCounter int

const returnRegister = "%eax"

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
		if addrDesc[_memloc].memLocation == "" {
			Code += fmt.Sprintf("movl %s, %s\n", addrDesc[_memloc].regLocation, reg)
			delete(regDesc[addrDesc[_memloc].regLocation], _memloc)
		} else {
			Code += fmt.Sprintf("movl (%s), %s\n", _memloc.MemoryLocation, reg)
		}
		regDesc[reg][_memloc] = true
		addrDesc[_memloc] = address{
			regLocation: reg,
			memLocation: addrDesc[_memloc].memLocation,
		}
	} else {
		Code += fmt.Sprintf("movl %s, %s\n", memloc.(*SymbolTableLiteralEntry).SymbolTableString(), reg)
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
	case SETRET:
		regDesc[returnRegister][ins.Arg1.(*SymbolTableVariableEntry)] = true
		addrDesc[ins.Arg1.(*SymbolTableVariableEntry)] = address{
			regLocation: returnRegister,
			memLocation: "",
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

func updateVariable(variable *SymbolTableVariableEntry, register MachineRegister) {
	if curreg := addrDesc[variable].regLocation; curreg != "" {
		delete(regDesc[curreg], variable)
	}
	regDesc[register][variable] = true
	addrDesc[variable] = address{
		regLocation: register,
		memLocation: "",
	}
}

func genUOpCode(ins IRIns, regs [3]registerResult) {
	spill(regs[2].Spill)
	switch ins.Op {
	case NEG:
		var valueString string
		if regs[0].Register == "" {
			valueString = ins.Arg1.SymbolTableString()
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

	dst := ins.Dst.(*SymbolTableVariableEntry)
	updateVariable(dst, regs[2].Register)
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

func genSOpCode(ins IRIns, regs [3]registerResult) {
	load(regs[0], ins.Arg1)
	spill(regs[2].Spill)
	Code += fmt.Sprintf("movl %s, %s\n", regs[0].Register, regs[2].Register)

	var valueString string
	if regs[1].Register == "" {
		valueString = ins.Arg2.SymbolTableString()
	} else {
		load(regs[1], ins.Arg2)
		valueString = "%cl"
	}

	switch ins.Op {
	case BSL:
		Code += fmt.Sprintf("shl %s, %s\n", valueString, regs[2].Register)
	case BSR:
		Code += fmt.Sprintf("shr %s, %s\n", valueString, regs[2].Register)
	}

	dst := ins.Dst.(*SymbolTableVariableEntry)
	updateVariable(dst, regs[2].Register)
}

func genBOpCode(ins IRIns, regs [3]registerResult) {
	load(regs[0], ins.Arg1)
	spill(regs[2].Spill)
	Code += fmt.Sprintf("movl %s, %s\n", regs[0].Register, regs[2].Register)

	var valueString string
	if regs[1].Register == "" {
		valueString = ins.Arg2.SymbolTableString()
	} else {
		load(regs[1], ins.Arg2)
		valueString = string(regs[1].Register)
	}

	switch ins.Op {
	case ADD:
		Code += fmt.Sprintf("addl %s, %s\n", valueString, regs[2].Register)
	case SUB:
		Code += fmt.Sprintf("subl %s, %s\n", valueString, regs[2].Register)
	case MUL:
		Code += fmt.Sprintf("imul %s, %s\n", valueString, regs[2].Register)
	case AND:
		fallthrough
	case BAND:
		Code += fmt.Sprintf("andl %s, %s\n", valueString, regs[2].Register)
	case OR:
		fallthrough
	case BOR:
		Code += fmt.Sprintf("orl %s, %s\n", valueString, regs[2].Register)
	case XOR:
		Code += fmt.Sprintf("xorl %s, %s\n", valueString, regs[2].Register)
	}

	dst := ins.Dst.(*SymbolTableVariableEntry)
	updateVariable(dst, regs[2].Register)
}

func genLOpCode(ins IRIns, regs [3]registerResult) {
	load(regs[0], ins.Arg1)
	spill(regs[2].Spill)

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
	}

	Code += fmt.Sprintf("movl $0, %s\n", regs[2].Register)
	Code += fmt.Sprintf("jmp _logic_end_%d", logicCounter)
	Code += fmt.Sprintf("_logic_start_%d:", logicCounter)
	Code += fmt.Sprintf("movl $1, %s\n", regs[2].Register)
	Code += fmt.Sprintf("_logic_end_%d:", logicCounter)
	logicCounter++

	dst := ins.Dst.(*SymbolTableVariableEntry)
	updateVariable(dst, regs[2].Register)
}

func genDOpCode(ins IRIns, regs [3]registerResult) {

	load(regs[0], ins.Arg1)
	load(regs[1], ins.Arg2)
	spill(regDesc["%edx"])

	Code += "movl $0, %edx\n"
	Code += fmt.Sprintf("idiv %s\n", regs[1].Register)

	dst := ins.Dst.(*SymbolTableVariableEntry)
	updateVariable(dst, regs[2].Register)
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
	case LOP:
		genLOpCode(ins, regs)
	case SOP:
		genSOpCode(ins, regs)
	case UOP:
		genUOpCode(ins, regs)
	case CBR:
		genCBRCode(ins, regs)
	case BOP:
		genBOpCode(ins, regs)
	case DOP:
		genDOpCode(ins, regs)
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
			if ins.Typ == LBL && i == 0 {
				target := ins.Dst.(*SymbolTableTargetEntry).Target
				if strings.HasPrefix(target, "_func") {
					Code += "push %ebp\n"
					Code += "movl %esp, %ebp\n"
				}
			}

			if i == len(bbl.Block)-1 && !shouldPreSave(ins) {
				saveBBL()
			}

		}
		clearBBL()
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

func clearBBL() {
	for _, variables := range regDesc {
		for variable := range variables {
			addrDesc[variable] = address{
				regLocation: "",
				memLocation: variable.MemoryLocation,
			}
		}
	}
	for register := range regDesc {
		regDesc[register] = make(map[*SymbolTableVariableEntry]bool)
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

main:
call _func_main
call exit

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
