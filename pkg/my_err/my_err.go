package my_err

import (
	"errors"
)

var (
	ErrUserExists   = errors.New("User already exists")
	ErrUserNotFound = errors.New("User not found")
	ErrNoAuth       = errors.New("No authentication session found")

	ErrEmptyTitle   = errors.New("Task title cannot be empty")
	ErrTaskNotFound = errors.New("User does not have task with given ID")
)
