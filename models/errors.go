package models

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("your requested item is not found")
	ErrBadParamInput       = errors.New("given param is not valid")
)
