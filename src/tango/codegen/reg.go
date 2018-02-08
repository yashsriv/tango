package codegen

import (
	"fmt"
	"math"
)

type address struct {
	regLocation string
	memLocation string
}

// Initialization of registers
var registers = map[string]*SymbolTableRegisterEntry{
	"r1": nil,
	"r2": nil,
	"r3": nil,
	"r4": nil,
}

//Initialization of variables
var variables = make(map[*SymbolTableRegisterEntry]address)

func assignRegister(i *SymbolTableRegisterEntry, uinfo map[*SymbolTableRegisterEntry]UseInfo,
	cannotbeReplaced map[string]bool, dst *SymbolTableRegisterEntry, canReplace bool) {

	if val, ok := variables[i]; ok && val.regLocation != "" {
		cannotbeReplaced[val.regLocation] = true
		return
	}

	for key, value := range registers {
		if value == nil {
			registers[key] = i
			// Update current address of i
			currentAddress := variables[i]
			currentAddress.regLocation = key
			variables[i] = currentAddress

			// If value is true it is not used and maybe replaced by dst if required
			cannotbeReplaced[key] = uinfo[i].NextUse == -1
			return
		}
	}

	score := make(map[string]int)

	for key, value := range registers {
		if val, ok := cannotbeReplaced[key]; ok && !(dst == i && val) {
			score[key] = math.MaxInt32
		}
		if variables[value].memLocation == "" {
			if !(canReplace && value == dst) {
				if uinfo[value].NextUse != -1 {
					score[key]++
				}
			}
		}
	}

	minScore := math.MaxInt32
	minScoreReg := ""

	for register, s := range score {
		if s < minScore {
			minScore = s
			minScoreReg = register
		}
	}

	// TODO: Return the fact that spilling needs to be done for this register
	// TODO: return the register being returned
	registers[minScoreReg] = i
	currentAddress := variables[i]
	currentAddress.regLocation = minScoreReg
	variables[i] = currentAddress

	cannotbeReplaced[minScoreReg] = uinfo[i].NextUse == -1
	return
}

func getReg(ins IRIns, uinfo map[*SymbolTableRegisterEntry]UseInfo) {

	instructionType := ins.Typ

	cannotBeReplaced := make(map[string]bool)

	switch instructionType {
	case BOP:
		canReplace := ins.Dst != ins.Arg1 && ins.Dst != ins.Arg2
		dst := ins.Dst.(*SymbolTableRegisterEntry)
		if i, isRegister := ins.Arg1.(*SymbolTableRegisterEntry); isRegister {
			//  i is a SymbolTableRegister
			assignRegister(i, uinfo, cannotBeReplaced, dst, canReplace)
		}
		if i, isRegister := ins.Arg2.(*SymbolTableRegisterEntry); isRegister {
			//  i is a SymbolTableRegister
			assignRegister(i, uinfo, cannotBeReplaced, dst, canReplace)
		}
		assignRegister(dst, uinfo, cannotBeReplaced, dst, canReplace)
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

}
