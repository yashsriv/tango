package ast

import (
	"fmt"
	"tango/src/codegen"
)

func ConstDecl(declnamelist, types, exprlist Attrib) (*AddrCode, error) {
	declnamelistAs, ok := declnamelist.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", declnamelist)
	}
	exprlistAs, ok := exprlist.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", exprlist)
	}

	if len(declnamelistAs) != len(exprlistAs) {
		return nil, fmt.Errorf("unequal number of elements in lhs and rhs: %d and %d", len(declnamelistAs), len(exprlistAs))
	}

	code := make([]codegen.IRIns, 0)

	for i := range declnamelistAs {
		entry, ok := declnamelistAs[i].Symbol.(*codegen.SymbolTableVariableEntry)
		if !ok {
			return nil, fmt.Errorf("lhs %s of expression should be a literal", declnamelistAs[i])
		}
		if entry.Assignments != 0 {
			return nil, fmt.Errorf("const %s can only be assigned once", entry.SymbolTableString())
		}
		entry.Assignments++
		if entry.Declared {
			return nil, fmt.Errorf("%s is being declared twice in this scope", entry.SymbolTableString())
		}
		entry.Declared = true
		code = append(code, exprlistAs[i].Code...)
		ins := codegen.IRIns{
			Typ:  codegen.ASN,
			Op:   codegen.ASNO,
			Dst:  declnamelistAs[i].Symbol,
			Arg1: exprlistAs[i].Symbol,
		}
		code = append(code, ins)
	}

	addrcode := &AddrCode{
		Code: code,
	}

	return addrcode, nil
}

func MultConstDecl(decl, decllist Attrib) (*AddrCode, error) {
	mergedList, err := MergeCodeList(decllist)
	if err != nil {
		return nil, err
	}
	declAsAddr, ok := decl.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", decl)
	}
	code := append(declAsAddr.Code, mergedList.Code...)
	addrcode := &AddrCode{
		Code: code,
	}

	return addrcode, nil
}
