package errors

import "errors"

var (
	ErrUnsupportedGrantType = errors.New("unsupported grant_type")
	ErrInvalidClient        = errors.New("invalid client")
	ErrInternalServer       = errors.New("internal server issue")
)
