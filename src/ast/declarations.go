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
