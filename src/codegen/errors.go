package codegen

import "errors"

// ErrAlreadyExists is error when symbol already exists in table
var ErrAlreadyExists = errors.New("symbol already exists in table")

// ErrDoesntExist is error when symbol is not in table
var ErrDoesntExist = errors.New("symbol doesn't exist in table")

// ErrEmptyTableStack is an error thrown when trying to pop off empty table stack
var ErrEmptyTableStack = errors.New("expected tableStack to never be empty")
