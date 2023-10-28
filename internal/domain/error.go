package domain

import "github.com/pkg/errors"

// all application domain level errors stored here
var (
	ErrNotFound        = errors.New("not found")
	ErrUserExists      = errors.New("user exists")
	ErrConversion      = errors.New("conversion error")
	ErrConfig          = errors.New("config error")
	ErrPasswordWrong   = errors.New("password wrong")
	ErrRePasswordWrong = errors.New("password and repeat password not equal")
)
