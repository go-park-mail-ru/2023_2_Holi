package domain

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("requested Item is not found")
	ErrBadRequest          = errors.New("request is not valid")
	ErrUnauthorized        = errors.New("need to authorize")
	ErrWrongCredentials    = errors.New("username or password is invalid")
	ErrAlreadyExists       = errors.New("resource already exists")
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case ErrAlreadyExists:
		return http.StatusForbidden
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrWrongCredentials:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
