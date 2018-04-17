package ast

import (
	"fmt"
	"tango/src/codegen"
)

// NewList creates a new list of attrib with initial element
func NewList(el Attrib) ([]*AddrCode, error) {
	list := make([]*AddrCode, 0)
	if el != nil {
		elAsAddrCode, ok := el.(*AddrCode)
		if !ok {
			return nil, fmt.Errorf("[NewList] unable to type cast %v to *AddrCode", el)
		}
		list = append(list, elAsAddrCode)
	}
	return list, nil
}

// AddToList adds an element to the list
func AddToList(list Attrib, el Attrib) ([]*AddrCode, error) {
	asList, ok := list.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[AddToList] unable to type cast %v to []*AddrCode", list)
	}
	elAsAddrCode, ok := el.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[AddToList] unable to type cast %v to *AddrCode", el)
	}
	return append(asList, elAsAddrCode), nil
}

// MergeCodeList can merge a list of codes
func MergeCodeList(list Attrib) (*AddrCode, error) {
	asList, ok := list.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[MergeCodeList] unable to typecast %v to []*AddrCode", list)
	}
	code := make([]codegen.IRIns, 0)
	for _, v := range asList {
		code = append(code, v.Code...)
	}
	addrcode := &AddrCode{
		Code: code,
	}
	return addrcode, nil
}

// NewIdList creates a new list of identifiers with initial element
func NewIdList(el Attrib) ([]*codegen.VariableEntry, error) {
	list := make([]*codegen.VariableEntry, 0)
	if el != nil {
		elAsToken, ok := el.(*codegen.VariableEntry)
		if !ok {
			return nil, fmt.Errorf("[NewIdList] unable to type cast %v to *SymbolTableVariableEntry", el)
		}
		list = append(list, elAsToken)
	}
	return list, nil
}

// AddToIdList adds an element to the list
func AddToIdList(list Attrib, el Attrib) ([]*codegen.VariableEntry, error) {
	asList, ok := list.([]*codegen.VariableEntry)
	if !ok {
		return nil, fmt.Errorf("[AddToIdList] unable to type cast %v to []*SymbolTableVariableEntry", list)
	}
	elAsAddrCode, ok := el.(*codegen.VariableEntry)
	if !ok {
		return nil, fmt.Errorf("[AddToIdList] unable to type cast %v to *SymbolTableVariableEntry", el)
	}
	return append(asList, elAsAddrCode), nil
}

// NewIfElseList creates a new list of attrib with initial element
func NewIfElseList(el Attrib) ([]*ifElse, error) {
	list := make([]*ifElse, 0)
	if el != nil {
		elAsAddrCode, ok := el.(*ifElse)
		if !ok {
			return nil, fmt.Errorf("[NewIfElseList] unable to type cast %v to *ifElse", el)
		}
		list = append(list, elAsAddrCode)
	}
	return list, nil
}

// AddToIfElseList adds an element to the list
func AddToIfElseList(list Attrib, el Attrib) ([]*ifElse, error) {
	asList, ok := list.([]*ifElse)
	if !ok {
		return nil, fmt.Errorf("[AddToIfElseList] unable to type cast %v to []*ifElse", list)
	}
	elAsAddrCode, ok := el.(*ifElse)
	if !ok {
		return nil, fmt.Errorf("[AddToIfElseList] unable to type cast %v to *ifElse", el)
	}
	return append(asList, elAsAddrCode), nil
}

// NewCaseBlockList creates a new list of attrib with initial element
func NewCaseBlockList() ([]*caseDecl, error) {
	list := make([]*caseDecl, 0)
	return list, nil
}

// AddToCaseBlockList adds an element to the list
func AddToCaseBlockList(list Attrib, el Attrib) ([]*caseDecl, error) {
	asList, ok := list.([]*caseDecl)
	if !ok {
		return nil, fmt.Errorf("[AddToCaseBlockList] unable to type cast %v to []*caseDecl", list)
	}
	elAsAddrCode, ok := el.(*caseDecl)
	if !ok {
		return nil, fmt.Errorf("[AddToCaseBlockList] unable to type cast %v to *caseDecl", el)
	}
	return append(asList, elAsAddrCode), nil
}
