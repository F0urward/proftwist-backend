package errs

import (
	"errors"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrBusinessLogic      = errors.New("business logic error")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInternal           = errors.New("internal server error")
)

func IsNotFoundError(err error) bool {
	return err != nil && (err.Error() == ErrNotFound.Error() ||
		contains(err.Error(), "not found") ||
		contains(err.Error(), "NotFound"))
}

func IsBusinessLogicError(err error) bool {
	return err != nil && (err.Error() == ErrBusinessLogic.Error() ||
		contains(err.Error(), "business logic") ||
		contains(err.Error(), "invalid") ||
		contains(err.Error(), "validation"))
}

func IsAlreadyExistsError(err error) bool {
	return err != nil && (err.Error() == ErrAlreadyExists.Error() ||
		contains(err.Error(), "already exists"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			contains(s[1:], substr))))
}
