package errs

import (
	"errors"
)

var (
	ErrReadRequestData    = errors.New("failed to read request body")
	ErrParseRequestData   = errors.New("failed to parse request body")
	ErrNotFound           = errors.New("not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidID          = errors.New("invalid id format")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrBusinessLogic      = errors.New("business logic error")
	ErrProductNotApproved = errors.New("product not approved")
	ErrNotEnoughStock     = errors.New("not enough stock")

	ErrMissingToken           = errors.New("missing jwt token")
	ErrTokenRevoked           = errors.New("token revoked")
	ErrNoMetadata             = errors.New("metadata is not provided")
	ErrNoAuthHeader           = errors.New("authorization header is missing")
	ErrInvalidAuthFormat      = errors.New("invalid authorization header format")
	ErrInternal               = errors.New("internal server error")
	ErrInvalidProductPrice    = errors.New("invalid product price")
	ErrEmptyProductName       = errors.New("invalid product name")
	ErrInvalidProductQuantity = errors.New("invalid product quantity")
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
