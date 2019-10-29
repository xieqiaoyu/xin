package xin

import (
	"fmt"
)

//InternalError
type InternalError struct {
	Err error
}

//Error error interface
func (e *InternalError) Error() string {
	return e.Err.Error()
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

func NewInternalError(format string, a ...interface{}) error {
	return &InternalError{
		Err: fmt.Errorf(format, a...),
	}
}
