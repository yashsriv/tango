package ast

import (
	"fmt"
	"tango/src/codegen"
	"tango/src/token"
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
		return nil, fmt.Errorf("[AddToList] unable to type cast %v, %T to *AddrCode", el, el)
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

// NewIdentifierList creates a new list of identifiers with initial element
func NewIdentifierList(el Attrib) (map[string]bool, error) {
	list := make(map[string]bool, 0)
	if el != nil {
		elAsToken := string(el.(*token.Token).Lit)
		list[elAsToken] = true
	}
	return list, nil
}

// AddToIdentifierList adds an element to the list
func AddToIdentifierList(list Attrib, el Attrib) (map[string]bool, error) {
	asList, ok := list.(map[string]bool)
	if !ok {
		return nil, fmt.Errorf("[AddToIdList] unable to type cast %v to map[string]bool", list)
	}
	elAsToken := string(el.(*token.Token).Lit)
	_, ok = asList[elAsToken]
	if ok {
		return nil, fmt.Errorf("identifier being redefined")
	}
	asList[elAsToken] = true
	return asList, nil
}

func NewArgTypeList(el Attrib) ([]*ArgType, error) {
	list := make([]*ArgType, 0)
	if el != nil {
		elAsToken, ok := el.(*ArgType)
		if !ok {
			return nil, fmt.Errorf("[NewArgTypeList] unable to type cast %v to *ArgType", el)
		}
		list = append(list, elAsToken)
	}
	return list, nil
}

func AddToArgTypeList(list, el Attrib) ([]*ArgType, error) {
	asList, ok := list.([]*ArgType)
	if !ok {
		return nil, fmt.Errorf("[AddToArgTypeList] unable to type cast %v to []*ArgType", list)
	}
	elAsAddrCode, ok := el.(*ArgType)
	if !ok {
		return nil, fmt.Errorf("[AddToArgTypeList] unable to type cast %v to *ArgType", el)
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
