package domain

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
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

func GetHttpStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case ErrWrongCredentials:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized

	case ErrInvalidToken:
		return http.StatusBadRequest
	case ErrBadRequest:
		return http.StatusBadRequest

	case ErrNotFound:
		return http.StatusNotFound
	case ErrOutOfRange:
		return http.StatusNotFound

	case ErrAlreadyExists:
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}

func GetGrpcStatusCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	switch err {
	case ErrNotFound:
		return codes.NotFound

	default:
		return codes.Internal
	}
}
