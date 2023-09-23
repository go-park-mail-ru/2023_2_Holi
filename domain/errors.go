package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("your requested Item is not found")
	//ErrConflict = errors.New("your Item already exist")
	ErrBadParamInput = errors.New("given Param is not valid")
	ErrUnauthorized  = errors.New("you need to authorize")
)
