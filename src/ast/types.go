package ast

import (
	"errors"
	"fmt"
	"reflect"
	"tango/src/codegen"
)

// IsType checks if something is a type entry or not
func IsType(a Attrib) (codegen.TypeEntry, error) {
	if v, ok := a.(codegen.TypeEntry); ok {
		return v, nil
	}
	return nil, fmt.Errorf("%v should be a TypeEntry", a)
}

// EvalPtrType checks if something is a pointer
func EvalPtrType(a Attrib) (codegen.TypeEntry, error) {
	t, err := IsType(a)
	if err != nil {
		return nil, err
	}
	return codegen.PtrType{To: t}, nil
}

// EvalArrType Evaluates an array
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

// CheckOperandType checks types of operands
func CheckOperandType(op codegen.IROp, t codegen.TypeEntry) error {
	if opMap, ok := opTypeMap[op]; ok {
		if _, ok1 := opMap[t]; ok1 {
			return nil
		}
		return errors.New("not supported operation on this type")
	}
	return errors.New("no known type associated with this operation")
}

// Alloc allocates default values to variables
func Alloc(v *codegen.VariableEntry) []codegen.IRIns {
	code := make([]codegen.IRIns, 0)
	switch t := v.VType.(type) {
	case codegen.StructType:
		code = append(code,
			codegen.IRIns{
				Typ:  codegen.KEY,
				Op:   codegen.MALLOC,
				Arg1: &codegen.LiteralEntry{Value: len(t.FieldTypes)},
			}, codegen.IRIns{
				Typ:  codegen.KEY,
				Op:   codegen.SETRET,
				Arg1: v,
			})
		for index, val := range t.FieldTypes {
			entry := CreateTemporary(val)
			code = append(code, entry.Code...)
			code = append(code, Alloc(entry.Symbol.(*codegen.VariableEntry))...)
			code = append(code, codegen.IRIns{
				Typ:  codegen.KEY,
				Op:   codegen.PUT,
				Dst:  v,
				Arg1: &codegen.LiteralEntry{Value: index},
				Arg2: entry.Symbol,
			})
		}
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
	case *codegen.BasicType, codegen.PtrType:
		code = append(code, codegen.IRIns{
			Typ:  codegen.ASN,
			Op:   codegen.ASNO,
			Dst:  v,
			Arg1: &codegen.LiteralEntry{Value: 0},
		})
	}
	return code
}

// EvalCompType evaluates a complicated type
func EvalCompType(a, b Attrib) (*AddrCode, error) {
	switch t := a.(codegen.TypeEntry).(type) {
	case codegen.StructType:
		entry := CreateTemporary(t)
		code := entry.Code
		var fields []keyval
		if b != nil {
			fields = b.([]keyval)
		}
		code = append(code, Alloc(entry.Symbol.(*codegen.VariableEntry))...)
		for _, field := range fields {
			index, ok := t.FieldMap[field.key]
			if !ok {
				return nil, fmt.Errorf("unknown field %s in %v", field.key, t)
			}
			if !SameType(t.FieldTypes[index], field.val.Symbol.Type()) {
				return nil, fmt.Errorf("invalid type for field %s. expected %v, got %v", field.key, t.FieldTypes[index], field.val.Symbol.Type())
			}
			code = append(code, field.val.Code...)
			code = append(code, codegen.IRIns{
				Typ:  codegen.KEY,
				Op:   codegen.PUT,
				Dst:  entry.Symbol,
				Arg1: &codegen.LiteralEntry{Value: index},
				Arg2: field.val.Symbol,
			})
		}
		return &AddrCode{Symbol: entry.Symbol, Code: code}, nil
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

// SameType compares two types properly
func SameType(a1 codegen.TypeEntry, a2 codegen.TypeEntry) bool {
	return reflect.DeepEqual(a1, a2)
}

// EvalStructEmbed evaluates embedded structs
func EvalStructEmbed(a Attrib) (map[string]codegen.TypeEntry, error) {
	t, err := IsType(a)
	if err != nil {
		return nil, err
	}
	if s, ok := t.(codegen.StructType); ok {
		typeMap := make(map[string]codegen.TypeEntry)
		for key, index := range s.FieldMap {
			val := s.FieldTypes[index]
			typeMap[key] = val
		}
		return typeMap, nil
	}
	return nil, fmt.Errorf("can only embed struct types. found: %v", t)
}

// EvalStructIDList evaluates a list of identifiers
func EvalStructIDList(a, b Attrib) (map[string]codegen.TypeEntry, error) {
	t := b.(codegen.TypeEntry)
	fields := a.(map[string]bool)
	m := make(map[string]codegen.TypeEntry)
	for field := range fields {
		m[field] = t
	}
	return m, nil
}

// NewStructDeclList creates a new list of identifiers with initial element
func NewStructDeclList(a Attrib) (map[string]codegen.TypeEntry, error) {
	list := make(map[string]codegen.TypeEntry)
	if a != nil {
		aAsMap := a.(map[string]codegen.TypeEntry)
		list = aAsMap
	}
	return list, nil
}

// AddToStructDeclList adds an element to the list
func AddToStructDeclList(fieldMap Attrib, el Attrib) (map[string]codegen.TypeEntry, error) {
	asMap, ok := fieldMap.(map[string]codegen.TypeEntry)
	if !ok {
		return nil, fmt.Errorf("[AddToIdList] unable to type cast %v to map[string]codegen.TypeEntry", asMap)
	}
	elAsMap := el.(map[string]codegen.TypeEntry)
	for fields, types := range elAsMap {
		_, ok = asMap[fields]
		if ok {
			return nil, fmt.Errorf("identifier being redefined")
		}
		asMap[fields] = types
	}
	return asMap, nil
}

// EvalStructType evaluates a struct type
func EvalStructType(a Attrib) (codegen.TypeEntry, error) {
	var asMap map[string]codegen.TypeEntry
	if a != nil {
		var ok bool
		asMap, ok = a.(map[string]codegen.TypeEntry)
		if !ok {
			return nil, fmt.Errorf("[AddToIdList] unable to type cast %v to map[string]codegen.TypeEntry", asMap)
		}
	}
	fieldTypes := make([]codegen.TypeEntry, len(asMap))
	fieldMap := make(map[string]int)
	var index int
	for key, t := range asMap {
		fieldTypes[index] = t
		fieldMap[key] = index
		index++
	}
	return codegen.StructType{FieldTypes: fieldTypes, FieldMap: fieldMap}, nil
}
