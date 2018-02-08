package codegen

import "fmt"

type address struct {
	regLocation string
	memLocation string
}

// Initialization of registers
var registers = map[string]string{"r1": "", "r2": "", "r3": "", "r4": ""}

//Initialization of variables
var variables = make(map[string]address)

func assignRegister(i SymbolTableRegisterEntry) {

	if val, ok := variables[i.Register]; ok && val.regLocation != "" {
		return
	}

	for key, value := range registers {
		if value == "" {
			registers[key] = i.Register
			variables[i.Register] = address{
				regLocation: i.Register,
				memLocation: i.Register,
			}
			return
		}
	}

}

func getReg(ins IRIns) {

	instructionType := ins.Typ

	switch instructionType {
	case BOP:
		if i, isRegister := ins.Arg1.(SymbolTableRegisterEntry); isRegister {
			//  i is a SymbolTableRegister
			assignRegister(i)
		}
		if i, isRegister := ins.Arg2.(SymbolTableRegisterEntry); isRegister {
			//  i is a SymbolTableRegister
			assignRegister(i)
		}
		if i, isRegister := ins.Dst.(SymbolTableRegisterEntry); isRegister {
			//  i is a SymbolTableRegister
			assignRegister(i)
		}
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
