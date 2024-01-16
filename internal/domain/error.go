package domain

import "github.com/pkg/errors"

// all application domain level errors stored here
// TODO: make all errors for domain level
var (
	ErrNotFound          = errors.New("not found")
	ErrRegRequiredFields = errors.New("email and password are required")
	ErrNameRequired      = errors.New("name is required")
	ErrUserExists        = errors.New("user exists")
	ErrConversion        = errors.New("conversion error")
	ErrConfig            = errors.New("config error")
	ErrPasswordWrong     = errors.New("password wrong")
	ErrRePasswordWrong   = errors.New("password and repeat password not equal")
)
