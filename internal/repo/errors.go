package repo

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exists")
	ErrCannotCreate = errors.New("can not create")
	ErrCannotGet    = errors.New("can not get")
	ErrCannotDelete = errors.New("can not delete")
)
