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

func spill(entries map[*VariableEntry]bool) {
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

	if _memloc, isRegister := memloc.(*VariableEntry); isRegister {
		// If we are trying to shift a variable already in the register, ignore
		if _, isKey := regDesc[reg][_memloc]; isKey {
			return
		}
	}

	spill(regres.Spill)

	// Load the value onto the register
	// can be a virtual register or a constant
	if _memloc, isRegister := memloc.(*VariableEntry); isRegister {
		if addrDesc[_memloc].memLocation == nil {
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
		Code += fmt.Sprintf("movl $%d, %s\n", memloc.(*LiteralEntry).Value, reg)
	}
}

func updateVariable(variable *VariableEntry, register MachineRegister) {
	if curreg := addrDesc[variable].regLocation; curreg != "" {
		delete(regDesc[curreg], variable)
	}
	regDesc[register][variable] = true
	addrDesc[variable] = address{
		regLocation: register,
		memLocation: nil,
	}
}

func genDOpCode(ins IRIns, regs [3]registerResult) {

	load(regs[0], ins.Arg1)
	load(regs[1], ins.Arg2)
	spill(regDesc["%edx"])

	// We'll need to do sign extension
	Code += "cltd\n"
	Code += fmt.Sprintf("idiv %s\n", regs[1].Register)

	dst := ins.Dst.(*VariableEntry)
	updateVariable(dst, regs[2].Register)
}

func genOpCode(ins IRIns, regs [3]registerResult) {
	switch ins.Typ {
	case LBL:
		Code += fmt.Sprintf("%s:\n", ins.Dst.(*TargetEntry).Target)
	case JMP:
		Code += fmt.Sprintf("jmp %s\n", ins.Arg1.(*TargetEntry).Target)
	case KEY:
		genKeyCode(ins, regs)
	case ASN:
		load(regs[0], ins.Arg1)
		dst := ins.Dst.(*VariableEntry)
		updateVariable(dst, regs[0].Register)
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
	// TODO: handle missing variable in addrDesc gracefully
	for _, bbl := range BBLList {
		Code += "# Begin Basic Block\n"
		for i, ins := range bbl.Block {
			arg1res, arg2res, dstres := getReg(ins, bbl.Info[i])

			if i == len(bbl.Block)-1 && isEndBlock(ins.Typ, ins.Op) {
				saveBBL()
			}

			genOpCode(ins, [3]registerResult{arg1res, arg2res, dstres})
			if ins.Typ == LBL && i == 0 {
				target := ins.Dst.(*TargetEntry).Target
				if strings.HasPrefix(target, "_func") {
					Code += "push %ebp\n"
					Code += "movl %esp, %ebp\n"
				}
			}

			if i == len(bbl.Block)-1 && !isEndBlock(ins.Typ, ins.Op) {
				saveBBL()
			}

		}
		clearBBL()
		Code += "# End Basic Block\n"
	}
}

func genData() {
	Code += ".section .data\n\n"
	Code += "true:  .long 1\n"
	Code += "false: .long 0\n"
	Code += "_fmtint: .string \"%d\"\n"
	Code += "_fmtchar: .string \"%c\"\n"
	Code += "_fmtstr: .string \"%s\"\n"
	for _, symbol := range SymbolTable.symbolMap {
		switch v := symbol.(type) {
		case *VariableEntry:
			if v.Name != "true" && v.Name != "false" {
				Code += fmt.Sprintf("%s: .long 0\n", v.MemoryLocation)
			}
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
		regDesc[register] = make(map[*VariableEntry]bool)
	}
}

func saveBBL() {
	// Code += "\n\n# Saving Stuff\n"
	for register, variables := range regDesc {
		for variable := range variables {
			if addrDesc[variable].memLocation == nil {
				Code += fmt.Sprintf("movl %s, (%s)\n", register, variable.MemoryLocation)
			}
		}
	}
	// Code += "# Done Saving Stuff\n\n"
}

func genMisc() {
	Code += `.section .text
.globl main

main:
call _func_init
call _func_main
push $0
call exit

`
}

// GenerateASM generates the assemby code in Code
func GenerateASM() {
	genMisc()
	genCode()
	genData()
}
