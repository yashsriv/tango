package ast

import (
	"fmt"
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

// StringLit creates a new String Literal entry in the SymbolTable
func StringLit(l Attrib) (*AddrCode, error) {
	byteval := string(l.(*token.Token).Lit)
	location := fmt.Sprintf("string_%d", stringCounter)
	codegen.Strings[location] = byteval
	stringCounter++
	entry := &codegen.VariableEntry{
		MemoryLocation: codegen.GlobalMemory{Location: location},
		VType:          stringType,
		Constant:       true,
	}
	codegen.CreateAddrDescEntry(entry)
	addrcode := &AddrCode{
		Symbol: entry,
		Code:   nil,
	}
	return addrcode, nil
}

var stringCounter int
