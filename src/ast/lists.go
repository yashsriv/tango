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
			return nil, fmt.Errorf("unable to type cast %v to *AddrCode", el)
		}
		list = append(list, elAsAddrCode)
	}
	return list, nil
}

// AddToList adds an element to the list
func AddToList(list Attrib, el Attrib) ([]*AddrCode, error) {
	asList, ok := list.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to []*AddrCode", list)
	}
	elAsAddrCode, ok := el.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", el)
	}
	return append(asList, elAsAddrCode), nil
}

// MergeCodeList can merge a list of codes
func MergeCodeList(list Attrib) (*AddrCode, error) {
	asList, ok := list.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("unable to typecast %v to []*AddrCode", list)
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

// NewIfElseList creates a new list of attrib with initial element
func NewIfElseList(el Attrib) ([]*ifElse, error) {
	list := make([]*ifElse, 0)
	if el != nil {
		elAsAddrCode, ok := el.(*ifElse)
		if !ok {
			return nil, fmt.Errorf("unable to type cast %v to *AddrCode", el)
		}
		list = append(list, elAsAddrCode)
	}
	return list, nil
}

// AddToIfElseList adds an element to the list
func AddToIfElseList(list Attrib, el Attrib) ([]*ifElse, error) {
	asList, ok := list.([]*ifElse)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to []*AddrCode", list)
	}
	elAsAddrCode, ok := el.(*ifElse)
	if !ok {
		return nil, fmt.Errorf("unable to type cast %v to *AddrCode", el)
	}
	return append(asList, elAsAddrCode), nil
}
