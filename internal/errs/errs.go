package errs

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidShortLink = errors.New("invalid short link lenght (!= 10)")
	ErrDuplicate        = errors.New("duplicate key")

	ErrEmptyURL         = errors.New("url is empty")
	ErrInvalidURLFormat = errors.New("invalid url format")

	ErrInternal = errors.New("db error")
)
