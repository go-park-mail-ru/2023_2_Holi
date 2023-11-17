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
	ErrInvalidToken        = errors.New("session token is invalid")
	ErrAlreadyExists       = errors.New("resource already exists")
	ErrOutOfRange          = errors.New("id is out of range")
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case ErrAlreadyExists:
		return http.StatusForbidden
	case ErrWrongCredentials:
		return http.StatusForbidden

	case ErrInvalidToken:
		return http.StatusBadRequest
	case ErrBadRequest:
		return http.StatusBadRequest

	case ErrNotFound:
		return http.StatusNotFound
	case ErrOutOfRange:
		return http.StatusNotFound

	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrUnauthorized:
		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}
