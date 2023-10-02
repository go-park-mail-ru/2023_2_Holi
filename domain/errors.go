package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("requested Item is not found")
	ErrBadRequest          = errors.New("request is not valid")
	ErrUnauthorized        = errors.New("need to authorize")
	ErrWrongCredentials    = errors.New("username or password is invalid")
	ErrAlreadyExists       = errors.New("resource already exists")
)
