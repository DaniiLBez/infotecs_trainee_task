package repo

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exists")
)
