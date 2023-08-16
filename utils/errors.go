package utils

import "errors"

var (
	ErrRecordAlreadyExist = errors.New("record exist")
	ErrNotFound           = errors.New("not found")
	ErrInvalid            = errors.New("invalid")
)
