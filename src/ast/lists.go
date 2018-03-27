package ast

import "errors"

// NewList creates a new list of attrib with initial element
func NewList(el Attrib) ([]*AddrCode, error) {
	list := make([]*AddrCode, 1)
	elAsAddrCode, ok := el.(*AddrCode)
	if !ok {
		return nil, errors.New("Unable to type cast to AddrCode")
	}
	list[0] = elAsAddrCode
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
