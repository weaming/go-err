package goerr

import (
	"errors"
)

var (
	InternalError = NewErrorFromString("Internal error")
	IOError       = InternalError.Wrap(errors.New("I/O error"))
	DBError       = IOError.Wrap(errors.New("Database error"))
)
