package html

import "tango/src/ast"

// Entry is an entry in the html
type Entry struct {
	Prefix ast.Stack
	Suffix ast.Stack
	Node   *ast.Node
}
