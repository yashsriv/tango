package ast

import (
	"errors"
	"fmt"
	"reflect"
	"tango/src/codegen"
)

func IsType(a Attrib) (codegen.TypeEntry, error) {
	if v, ok := a.(codegen.TypeEntry); ok {
		return v, nil
	}
	return nil, fmt.Errorf("%v should be a TypeEntry", a)
}

func EvalPtrType(a Attrib) (codegen.TypeEntry, error) {
	t, err := IsType(a)
	if err != nil {
		return nil, err
	}
	return codegen.PtrType{To: t}, nil
}

func EvalArrType(a, b Attrib) (codegen.TypeEntry, error) {
	t, err := IsType(b)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return codegen.ArrType{Of: t, Size: 0}, nil
	}
	switch s := a.(*AddrCode).Symbol.(type) {
	case *codegen.LiteralEntry:
		if s.Value < 1 {
			return nil, fmt.Errorf("cannot use non-positive array size: %d", s.Value)
		}
		return codegen.ArrType{Of: t, Size: s.Value}, nil
	default:
		return nil, fmt.Errorf("non-constant array sizes are not supported")
	}
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
	codegen.NOT:   map[codegen.TypeEntry]bool{boolType: true},
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

func Alloc(v *codegen.VariableEntry) []codegen.IRIns {
	code := make([]codegen.IRIns, 0)
	switch t := v.VType.(type) {
	case codegen.ArrType:
		if t.Size != 0 {
			code = append(code,
				codegen.IRIns{
					Typ:  codegen.KEY,
					Op:   codegen.MALLOC,
					Arg1: &codegen.LiteralEntry{Value: t.Size},
				}, codegen.IRIns{
					Typ:  codegen.KEY,
					Op:   codegen.SETRET,
					Arg1: v,
				})
			for i := 0; i < t.Size; i++ {
				entry := CreateTemporary(t.Of)
				code = append(code, entry.Code...)
				code = append(code, Alloc(entry.Symbol.(*codegen.VariableEntry))...)
				code = append(code, codegen.IRIns{
					Typ:  codegen.KEY,
					Op:   codegen.PUT,
					Dst:  v,
					Arg1: &codegen.LiteralEntry{Value: i},
					Arg2: entry.Symbol,
				})
			}
		}
	case *codegen.BasicType:
		code = append(code, codegen.IRIns{
			Typ:  codegen.ASN,
			Op:   codegen.ASNO,
			Dst:  v,
			Arg1: &codegen.LiteralEntry{Value: 0},
		})
	}
	return code
}

func EvalCompType(a, b Attrib) (*AddrCode, error) {
	switch t := a.(codegen.TypeEntry).(type) {
	case codegen.ArrType:
		entry := CreateTemporary(t)
		code := entry.Code
		var expr []*AddrCode
		if b != nil {
			expr = b.([]*AddrCode)
		}
		if t.Size != 0 {
			if len(expr) != t.Size {
				return nil, fmt.Errorf("size of array is not as mentioned in type. expected %d, got %d", t.Size, len(expr))
			}
			code = append(code, Alloc(entry.Symbol.(*codegen.VariableEntry))...)
			for i := range expr {
				code = append(code, expr[i].Code...)
				code = append(code, codegen.IRIns{
					Typ:  codegen.KEY,
					Op:   codegen.PUT,
					Dst:  entry.Symbol,
					Arg1: &codegen.LiteralEntry{Value: i},
					Arg2: expr[i].Symbol,
				})
			}
		} else {
			entry.Symbol.(*codegen.VariableEntry).VType = codegen.ArrType{Of: t.Of, Size: len(expr)}
			code = append(code, Alloc(entry.Symbol.(*codegen.VariableEntry))...)
			entry.Symbol.(*codegen.VariableEntry).VType = t
			for i := range expr {
				code = append(code, expr[i].Code...)
				code = append(code, codegen.IRIns{
					Typ:  codegen.KEY,
					Op:   codegen.PUT,
					Dst:  entry.Symbol,
					Arg1: &codegen.LiteralEntry{Value: i},
					Arg2: expr[i].Symbol,
				})
			}
		}
		return &AddrCode{Symbol: entry.Symbol, Code: code}, nil
	default:
		return nil, ErrUnsupported
	}
}

func SameType(a1 codegen.TypeEntry, a2 codegen.TypeEntry) bool {
	return reflect.DeepEqual(a1, a2)
}
