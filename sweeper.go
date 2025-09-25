package sweeper

import "errors"

var (
	ErrBadCredentials = errors.New("email/password wrong combination")
	ErrorNotFound     = errors.New("error not found")
	ErrValidation     = errors.New("validation error")
)
