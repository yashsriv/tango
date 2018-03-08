package ast

// Attrib represents any generic element of the ast
type Attrib interface {
}

// Node represents a node
type Node struct {
	name string
}

func (n Node) String() string {
	return n.name
}

// Derivations are the all the derivations discovered in this parse
var Derivations map[*Node]Stack

func init() {
	Derivations = make(map[*Node]Stack)
}

// AddNode creates a node
func AddNode(name string, attribs ...Attrib) (Attrib, error) {
	node := &Node{
		name: name,
	}
	Derivations[node] = attribs
	return node, nil
}
