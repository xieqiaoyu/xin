package xin

import (
	"errors"
	"fmt"
)

type WrapError interface {
	error
	Wrap(error)
}

type E struct {
	Err error
}

func (e *E) Error() string {
	return e.Err.Error()
}

func (e *E) Unwrap() error {
	return e.Err
}

func (e *E) Wrap(err error) {
	e.Err = err
}

func WrapE(Err WrapError, format string, a ...interface{}) error {
	var wErr error
	if len(a) > 0 {
		wErr = fmt.Errorf(format, a...)
	} else {
		wErr = errors.New(format)
	}
	Err.Wrap(wErr)
	return Err
}

//InternalError
type InternalError struct{ E }
