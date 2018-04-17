package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
)

// FuncSign is the function's signature
func FuncSign(a, b, c, d Attrib) (*AddrCode, error) {
	// TODO: Handle other stuff like arg list, return type and method declarations
	identifier := string(a.(*token.Token).Lit)
	start := &codegen.TargetEntry{
		Target: fmt.Sprintf("_func_%s", identifier),
	}
	// Associating identifier with some entry
	err := codegen.SymbolTable.InsertSymbol(identifier, start)
	return &AddrCode{Symbol: start}, err
}
