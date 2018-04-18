package ast

import (
	"errors"
	"fmt"
	"tango/src/codegen"
)

// Attrib represents any generic element of the ast
type Attrib interface {
}

// Node represents a node
type Node struct {
	name string
}

func (n *Node) String() string {
	if n == nil {
		return ""
	}
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

// AddrCode is a struct representing the SymbolTableEntry and Code
type AddrCode struct {
	Symbol codegen.SymbolTableEntry
	Code   []codegen.IRIns
}

// ErrUnsupported is used to report unsupported errors in the code
var ErrUnsupported = errors.New("unsupported operation")

// NewSourceFile creates a source file from the decl list
func NewSourceFile(declList Attrib) (*AddrCode, error) {
	// Perform hoisting
	// TODO: get names and types of all functions available
	var initAddrCode *AddrCode
	initList := make([]*AddrCode, 0)
	funcList := make([]*AddrCode, 0)
	asList, ok := declList.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[NewSourceFile] unable to typecast %v to []*AddrCode", declList)
	}
	for _, v := range asList {
		if len(v.Code) == 0 {
			continue
		}
		if v.Code[0].Typ == codegen.LBL {
			// Function Declaration
			// We just need to check if it isn't the declaration for the
			// init function
			if v.Code[0].Dst.(*codegen.TargetEntry).Target == "_func_init" {
				initAddrCode = v
			} else {
				funcList = append(funcList, v)
			}
		} else {
			initList = append(initList, v)
		}
	}
	initCode, err := MergeCodeList(initList)
	if err != nil {
		return nil, err
	}
	funcCode, err := MergeCodeList(funcList)
	if err != nil {
		return nil, err
	}
	if initAddrCode != nil {
		initCode.Code = append(initCode.Code, initAddrCode.Code[1:]...)
	}
	code := make([]codegen.IRIns, 1)

	// Declare our own init function
	code[0] = codegen.IRIns{
		Typ: codegen.LBL,
		Dst: &codegen.TargetEntry{
			Target: "_func_init",
		},
	}
	code = append(code, initCode.Code...)
	code = append(code, codegen.IRIns{
		Typ: codegen.KEY,
		Op:  codegen.RET,
	})

	// Add remaining stuff
	code = append(code, funcCode.Code...)

	addrcode := &AddrCode{
		Code: code,
	}
	return addrcode, nil
}

var tempCount int

var predecID = []string{
	"bool", "byte", "error", "float32",
	"int", "int8", "int16", "int32", "rune", "string",
	"uint", "uint8", "uint16", "uint32", "uintptr",
}

var predecConst = []string{
	"true",
	"false",
}

var predecFunc = []string{
	"printf",
}

func init() {
	for _, v := range predecID {
		// TODO: Make this symboltable type entry or something
		codegen.SymbolTable.InsertSymbol(v, nil)
	}
	for _, v := range predecConst {
		codegen.SymbolTable.InsertSymbol(v, &codegen.VariableEntry{
			MemoryLocation: codegen.GlobalMemory{Location: v},
			Name:           v,
			Constant:       true,
		})
	}
	for _, v := range predecFunc {
		codegen.SymbolTable.InsertSymbol(v, &codegen.TargetEntry{
			Target: v,
		})
	}
}
