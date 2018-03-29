package ast

import (
	"errors"
	"fmt"
	"tango/src/codegen"
)

// NewList creates a new list of attrib with initial element
func NewList(el Attrib) ([]*AddrCode, error) {
	list := make([]*AddrCode, 0)
	if el != nil {
		elAsAddrCode, ok := el.(*AddrCode)
		if !ok {
			return nil, errors.New("Unable to type cast to AddrCode")
		}
		list = append(list, elAsAddrCode)
	}
	return list, nil
}

// AddToList adds an element to the list
func AddToList(list Attrib, el Attrib) ([]*AddrCode, error) {
	asList, ok := list.([]*AddrCode)
	if !ok {
		return nil, errors.New("Unable to type cast to list")
	}
	elAsAddrCode, ok := el.(*AddrCode)
	if !ok {
		return nil, errors.New("Unable to type cast to AddrCode")
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
