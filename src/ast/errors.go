package ast

import "errors"

var ErrShouldBeVariable = errors.New("identifier should be variable inorder to be used in this context")
