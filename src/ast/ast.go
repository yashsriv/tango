package ast

import (
	"fmt"
	"tango/src/token"
)

// Attrib represents any generic element of the ast
type Attrib interface {
}

type ourAttrib interface {
	// GenGraph(gographviz.Interface) string
	Name() string
	GenOutput()
}

type Stack []Attrib

func (s Stack) Empty() bool {
	return len(s) == 0
}

func (s Stack) Push(v Attrib) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, Attrib) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s Stack) String() string {
	str := ""
	for _, value := range s {
		switch v := value.(type) {
		case *token.Token:
			str += fmt.Sprintf("%q ", v)
		default:
			str += fmt.Sprintf("%s ", v)
		}
	}
	return str
}

// Node represents a node
type Node struct {
	name     string
	Children []ourAttrib
}

func (n Node) String() string {
	return n.name
}

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

func (n *Node) Name() string {
	return n.name
}

func (n *Node) GenOutput() {
	for _, val := range n.Children {
		fmt.Printf("%s => %s\n", n.Name(), val.Name())
		val.GenOutput()
	}
}
