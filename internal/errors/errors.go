package errors

import (
	"errors"
)

var (
	ErrInvalidOrderNumber            = errors.New("invalid order number")
	ErrOrderWasUploadedByCurrentUser = errors.New("the order was uploaded by current user")
	ErrOrderWasUploadedByAnotherUser = errors.New("the order was uploaded by another user")

	ErrUserLoginAlreadyExists       = errors.New("user login already exists")
	ErrInvalidUserLoginOrPassword   = errors.New("invalid login or password")
	ErrUserLoginAndPasswordRequired = errors.New("login and password are required")
	ErrUserUnauthorized             = errors.New("user unauthorized")

	ErrNegativeBalance = errors.New("negative balance")
)
