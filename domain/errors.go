package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("your requested Item is not found")
	ErrBadRequest          = errors.New("given Param is not valid")
	ErrUnauthorized        = errors.New("you need to authorize")
	ErrWrongCredentials    = errors.New("your username or password is invalid")
)
