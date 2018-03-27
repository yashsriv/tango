package ast

import (
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

// Identifier gets an identifier from the table
func Identifier(l Attrib) (*AddrCode, error) {
	identifier := string(l.(*token.Token).Lit)
	entry, err := codegen.InsertToSymbolTable("r" + identifier)
	if err != nil {
		return nil, err
	}
	addrcode := &AddrCode{
		Symbol: entry,
		Code:   nil,
	}
	return addrcode, nil
}
