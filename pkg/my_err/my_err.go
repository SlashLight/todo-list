package my_err

import (
	"errors"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrNoAuth       = errors.New("no authentication session found")

	ErrEmptyTitle   = errors.New("task title cannot be empty")
	ErrTaskNotFound = errors.New("user does not have task with given ID")

	ErrEmptyField = errors.New("field cannot be empty")
	ErrParseUUID  = errors.New("failed to parse UUID")
)
