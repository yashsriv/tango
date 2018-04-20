package ast

import (
	"fmt"
	"tango/src/token"
)

type keyval struct {
	key string
	val *AddrCode
}

// EvalKeyval evaluates a keyval
func EvalKeyval(a, b Attrib) (keyval, error) {
	identifier := string(a.(*token.Token).Lit)
	value := b.(*AddrCode)
	return keyval{identifier, value}, nil
}

// NewKeyvalList creates a new list of identifiers with initial element
func NewKeyvalList(el Attrib) ([]keyval, error) {
	list := make([]keyval, 0)
	if el != nil {
		elAsToken := el.(keyval)
		list = append(list, elAsToken)
	}
	return list, nil
}

// AddToKeyvalList adds an element to the list
func AddToKeyvalList(list Attrib, el Attrib) ([]keyval, error) {
	asList, ok := list.([]keyval)
	if !ok {
		return nil, fmt.Errorf("[AddToKeyvalList] unable to type cast %v to []keyval", list)
	}
	elAsToken := el.(keyval)
	return append(asList, elAsToken), nil
}
