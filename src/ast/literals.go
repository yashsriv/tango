package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

// IntLit creates a new Integer Literal entry in the SymbolTable
func IntLit(l Attrib) (*AddrCode, error) {
	byteval := string(l.(*token.Token).Lit)
	entry, err := codegen.InsertToSymbolTable("$" + byteval)
	if err != nil {
		return nil, err
	}
	addrcode := &AddrCode{
		Symbol: entry,
		Code:   nil,
	}
	return addrcode, nil
}

// Label gets a label from the table
func Label(l Attrib) (*AddrCode, error) {
	identifier := string(l.(*token.Token).Lit)
	if _, ok := codegen.AccSymbolMap(identifier); ok {
		return nil, fmt.Errorf("Identifier %s already used in this scope", identifier)
	}
	entry := &codegen.SymbolTableTargetEntry{
		Target: "_func_" + identifier,
	}
	codegen.InsertToSymbolMap(identifier, entry)
	addrcode := &AddrCode{
		Symbol: entry,
		Code:   nil,
	}
	return addrcode, nil
}
