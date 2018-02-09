package codegen

import (
	"fmt"
	"math"
)

type address struct {
	regLocation MachineRegister
	memLocation string
}

type registerResult struct {
	Register MachineRegister
	Spill    *SymbolTableRegisterEntry
}

// MachineRegister represents a register in a machine
type MachineRegister string

// Initialization of registers
var registers = map[MachineRegister]*SymbolTableRegisterEntry{
	"eax": nil,
	"ebx": nil,
	"ecx": nil,
	"edx": nil,
}

//Initialization of variables
var variables = make(map[*SymbolTableRegisterEntry]address)

func assignHelper(uinfo map[*SymbolTableRegisterEntry]UseInfo, dst *SymbolTableRegisterEntry, canReplace bool) func(*SymbolTableRegisterEntry) registerResult {

	// Create a closure storing the value of cannot be replaced

	cannotbeReplaced := make(map[MachineRegister]bool)

	return func(i *SymbolTableRegisterEntry) registerResult {
		// If variable to be assigned is already in a register
		// return that register
		if variables[i].regLocation != "" {
			cannotbeReplaced[variables[i].regLocation] = uinfo[i].NextUse == -1
			return registerResult{Register: variables[i].regLocation}
		}

		for register, variable := range registers {
			if variable == nil {

				if _, notReplace := cannotbeReplaced[register]; notReplace {
					continue
				}

				// If register doesn;t contain any variable, return that register

				// If value is true, it is not used anymore and maybe replaced by dst if required
				cannotbeReplaced[register] = uinfo[i].NextUse == -1
				return registerResult{Register: register}
			}
		}

		// Calculate scores of registers. We may still have to spill/not spill stuff
		score := make(map[MachineRegister]int)
		for key, value := range registers {
			// If the register has been allocated in this step to another variable, it should not be chosen
			// But, if we are choosing for the destination and the variable has no next user, we can choose this
			// so score shouldn't be math.MaxInt32 then
			if val, ok := cannotbeReplaced[key]; ok && !(dst == i && val) {
				score[key] = math.MaxInt32
			} else if variables[value].memLocation == "" {
				// variables which are not in any memLocation
				// we can overwrite the dst register if it is not going to be used
				// in the future. So score shouldn't be incremented
				// NOTE: Issue here. Variable isn't stored in memory and we are replacing it.
				// If variable is used in another block, we have lost the value
				if !(canReplace && value == dst) {
					if uinfo[value].NextUse != -1 {
						// If the variable we are replacing has to be used in the future,
						// add 1 to score.
						score[key]++
					}
				}
			}
		}

		minScore := math.MaxInt32
		var minScoreReg MachineRegister

		for register, s := range score {
			if s < minScore {
				minScore = s
				minScoreReg = register
			}
		}

		cannotbeReplaced[minScoreReg] = uinfo[i].NextUse == -1

		if minScore != 0 {
			return registerResult{Register: minScoreReg, Spill: registers[minScoreReg]}
		}

		return registerResult{Register: minScoreReg}
	}

}

func getReg(ins IRIns, uinfo map[*SymbolTableRegisterEntry]UseInfo) (arg1res, arg2res, dstres registerResult) {

	instructionType := ins.Typ

	switch instructionType {
	case BOP:
		canReplace := ins.Dst != ins.Arg1 && ins.Dst != ins.Arg2
		dst := ins.Dst.(*SymbolTableRegisterEntry)
		assignRegister := assignHelper(uinfo, dst, canReplace)
		if arg1, isRegister := ins.Arg1.(*SymbolTableRegisterEntry); isRegister {
			//  i is a SymbolTableRegister
			arg1res = assignRegister(arg1)
		}
		if arg2, isRegister := ins.Arg2.(*SymbolTableRegisterEntry); isRegister {
			//  i is a SymbolTableRegister
			arg2res = assignRegister(arg2)
		}
		dstres = assignRegister(dst)
	case UOP:
		fmt.Println("UOP")
	case CBR:
		fmt.Println("CBR")
	case JMP:
		fmt.Println("JMP")
	case ASN:
		fmt.Println("ASN")
	case KEY:
		fmt.Println("KEY")
	case INV:
		fmt.Println("INV")
	}
	return
}
