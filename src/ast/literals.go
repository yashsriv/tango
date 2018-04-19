package ast

import (
	"strconv"
	"tango/src/codegen"
	"tango/src/token"
)

// IntLit creates a new Integer Literal entry in the SymbolTable
func IntLit(l Attrib) (*AddrCode, error) {
	byteval := string(l.(*token.Token).Lit)
	val, err := strconv.ParseInt(byteval, 0, 32)
	if err != nil {
		return nil, err
	}
	entry := &codegen.LiteralEntry{
		Value: int(val),
		LType: intType,
	}
	addrcode := &AddrCode{
		Symbol: entry,
		Code:   nil,
	}
	return addrcode, nil
}
