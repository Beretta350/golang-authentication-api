package errs

import "errors"

var (
	ErrUsernameOrPasswordMismatch error = errors.New("invalid username or password")
	ErrMissingDataInRequest       error = errors.New("missing data in request")
)
