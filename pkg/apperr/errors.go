package apperr

import "errors"

var (
	ErrValidation    = errors.New("validation_error")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrNotFound      = errors.New("not_found")
	ErrConflict      = errors.New("conflict")
	ErrUnprocessable = errors.New("unprocessable_entity")
	ErrInternal      = errors.New("internal_error")
)

type WithDetails struct {
	Kind    error
	Message string
	Details interface{}
}

func (e *WithDetails) Error() string { return e.Message }
func (e *WithDetails) Unwrap() error { return e.Kind }

func Wrap(kind error, message string) error {
	return &WithDetails{Kind: kind, Message: message}
}

func WrapWithDetails(kind error, message string, details interface{}) error {
	return &WithDetails{Kind: kind, Message: message, Details: details}
}
