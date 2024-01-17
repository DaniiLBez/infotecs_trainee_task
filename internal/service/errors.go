package service

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exist")
	ErrCannotCreateUser  = errors.New("can not create user")
	ErrUserNotFound      = errors.New("user not found")
	ErrCannotGetUser     = errors.New("can not get user")

	ErrCannotSignToken  = errors.New("cannot sign token")
	ErrCannotParseToken = errors.New("cannot parse token")
)
