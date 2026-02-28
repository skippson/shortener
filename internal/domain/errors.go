package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExist     = errors.New("already exists")
	ErrInvalidURL       = errors.New("invalid url")
	ErrInvalidShortened = errors.New("invalid shortened")
)
