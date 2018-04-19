package ast

import (
	"errors"
	"fmt"
	"tango/src/codegen"
)

func IsType(a Attrib) (codegen.TypeEntry, error) {
	if v, ok := a.(codegen.TypeEntry); ok {
		return v, nil
	}
	return nil, fmt.Errorf("%v should be a TypeEntry", a)
}

var opTypeMap = map[codegen.IROp]map[codegen.TypeEntry]bool{
	codegen.ADD:   map[codegen.TypeEntry]bool{intType: true},
	codegen.SUB:   map[codegen.TypeEntry]bool{intType: true},
	codegen.MUL:   map[codegen.TypeEntry]bool{intType: true},
	codegen.BAND:  map[codegen.TypeEntry]bool{intType: true},
	codegen.BOR:   map[codegen.TypeEntry]bool{intType: true},
	codegen.BNOT:  map[codegen.TypeEntry]bool{intType: true},
	codegen.NEG:   map[codegen.TypeEntry]bool{intType: true},
	codegen.XOR:   map[codegen.TypeEntry]bool{intType: true},
	codegen.DIV:   map[codegen.TypeEntry]bool{intType: true},
	codegen.REM:   map[codegen.TypeEntry]bool{intType: true},
	codegen.BSL:   map[codegen.TypeEntry]bool{intType: true},
	codegen.BSR:   map[codegen.TypeEntry]bool{intType: true},
	codegen.BRGT:  map[codegen.TypeEntry]bool{intType: true},
	codegen.BRGTE: map[codegen.TypeEntry]bool{intType: true},
	codegen.BRLT:  map[codegen.TypeEntry]bool{intType: true},
	codegen.BRLTE: map[codegen.TypeEntry]bool{intType: true},
	codegen.INC:   map[codegen.TypeEntry]bool{intType: true},
	codegen.DEC:   map[codegen.TypeEntry]bool{intType: true},
	codegen.AND:   map[codegen.TypeEntry]bool{boolType: true},
	codegen.OR:    map[codegen.TypeEntry]bool{boolType: true},
	codegen.BREQ:  map[codegen.TypeEntry]bool{intType: true, boolType: true},
	codegen.BRNEQ: map[codegen.TypeEntry]bool{intType: true, boolType: true},
}

func CheckOperandType(op codegen.IROp, t codegen.TypeEntry) error {
	if opMap, ok := opTypeMap[op]; ok {
		if _, ok1 := opMap[t]; ok1 {
			return nil
		}
		return errors.New("not supported operation on this type")
	}
	return errors.New("no known type associated with this operation")
}
