package errs

import "errors"

var (
	ErrUsernameOrPasswordMismatch error = errors.New("invalid username or password")
)
