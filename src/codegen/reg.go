package codegen

import (
	"log"
	"math"
)

type address struct {
	regLocation MachineRegister
	memLocation string
}

type registerResult struct {
	Register MachineRegister
	Spill    map[*SymbolTableVariableEntry]bool
}

// MachineRegister represents a register in a machine
type MachineRegister string

// Initialization of regDesc
var regDesc = map[MachineRegister]map[*SymbolTableVariableEntry]bool{
	"%eax": map[*SymbolTableVariableEntry]bool{},
	"%ebx": map[*SymbolTableVariableEntry]bool{},
	"%ecx": map[*SymbolTableVariableEntry]bool{},
	"%edx": map[*SymbolTableVariableEntry]bool{},
}

//Initialization of addrDesc
var addrDesc = make(map[*SymbolTableVariableEntry]address)

func assignHelper(uinfo map[*SymbolTableVariableEntry]UseInfo, dst *SymbolTableVariableEntry, canReplace bool) func(*SymbolTableVariableEntry) registerResult {

	// Create a closure storing the value of cannot be replaced

	cannotbeReplaced := make(map[MachineRegister]bool)

	return func(i *SymbolTableVariableEntry) registerResult {
		// If variable to be assigned is already in a register
		// return that register
		if addrDesc[i].regLocation != "" {
			cannotbeReplaced[addrDesc[i].regLocation] = uinfo[i].NextUse == -1
			return registerResult{Register: addrDesc[i].regLocation}
		}

		for register, variable := range regDesc {
			if len(variable) == 0 {

				if _, notReplace := cannotbeReplaced[register]; notReplace {
					continue
				}

				// If register doesn;t contain any variable, return that register

				// If value is true, it is not used anymore and maybe replaced by dst if required
				cannotbeReplaced[register] = uinfo[i].NextUse == -1
				return registerResult{Register: register}
			}
		}

		// Calculate scores of regDesc. We may still have to spill/not spill stuff
		score := make(map[MachineRegister]int)
		for key, values := range regDesc {
			// If the register has been allocated in this step to another variable, it should not be chosen
			// But, if we are choosing for the destination and the variable has no next user, we can choose this
			// so score shouldn't be math.MaxInt32 then
			if val, ok := cannotbeReplaced[key]; ok && !(dst == i && val) {
				score[key] = math.MaxInt32
			} else {
				for value := range values {
					if addrDesc[value].memLocation == "" {
						// addrDesc which are not in any memLocation
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
			return registerResult{Register: minScoreReg, Spill: regDesc[minScoreReg]}
		}

		return registerResult{Register: minScoreReg}
	}

}

// getReg returns an allocation of regDesc for the operands
func getReg(ins IRIns, uinfo map[*SymbolTableVariableEntry]UseInfo) (arg1res, arg2res, dstres registerResult) {

	instructionType := ins.Typ

	switch instructionType {
	case BOP:
		canReplace := ins.Dst != ins.Arg1 && ins.Dst != ins.Arg2
		dst := ins.Dst.(*SymbolTableVariableEntry)
		assignRegister := assignHelper(uinfo, dst, canReplace)
		if arg1, isRegister := ins.Arg1.(*SymbolTableVariableEntry); isRegister {
			//  i is a SymbolTableRegister
			arg1res = assignRegister(arg1)
		}
		if arg2, isRegister := ins.Arg2.(*SymbolTableVariableEntry); isRegister {
			//  i is a SymbolTableRegister
			arg2res = assignRegister(arg2)
		}
		dstres = assignRegister(dst)
	case UOP:
		canReplace := ins.Dst != ins.Arg1
		dst := ins.Dst.(*SymbolTableVariableEntry)
		assignRegister := assignHelper(uinfo, dst, canReplace)
		if arg1, isRegister := ins.Arg1.(*SymbolTableVariableEntry); isRegister {
			//  i is a SymbolTableRegister
			arg1res = assignRegister(arg1)
		}
		dstres = assignRegister(dst)
	case CBR:
		assignRegister := assignHelper(uinfo, nil, false)
		if arg1, isRegister := ins.Arg1.(*SymbolTableVariableEntry); isRegister {
			//  i is a SymbolTableRegister
			arg1res = assignRegister(arg1)
		}
		if arg2, isRegister := ins.Arg2.(*SymbolTableVariableEntry); isRegister {
			//  i is a SymbolTableRegister
			arg2res = assignRegister(arg2)
		}
	case JMP:
		// DO Nothing
	case LBL:
		// Do Nothing
	case ASN:
		canReplace := ins.Dst != ins.Arg1
		dst := ins.Dst.(*SymbolTableVariableEntry)
		assignRegister := assignHelper(uinfo, dst, canReplace)
		arg1res = assignRegister(dst)
		dstres = arg1res
	case KEY:
		assignRegister := assignHelper(uinfo, nil, false)
		if !(ins.Op == RET || ins.Op == HALT) {
			if arg1, isRegister := ins.Arg1.(*SymbolTableVariableEntry); isRegister {
				//  i is a SymbolTableRegister
				arg1res = assignRegister(arg1)
			}
		}
	case INV:
		log.Fatalf("Invalid Instruction found. Aborting!!\n")
	}
	return
}
