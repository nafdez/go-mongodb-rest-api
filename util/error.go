package util

import (
	"errors"
)

var (
	ErrNoUsernameOrPasswordProvided = errors.New("no username or password provided")
	ErrInvalidUsernameOrPassword    = errors.New("username or password are wrong")
	ErrNoValidTokenProvided         = errors.New("no valid token provided")
	ErrUserNotFound                 = errors.New("user not found")
	ErrUserAlreadyExists            = errors.New("user already exist")
)
