package xin

import (
	"errors"
	"fmt"
	"runtime"
)

type WrapError interface {
	error
	Wrap(error)
}

type E struct {
	File string
	Line int
	Err  error
}

func (e *E) Error() string {
	return fmt.Sprintf("%s (%s:%d)", e.Err, e.File, e.Line)
}

func (e *E) Unwrap() error {
	return e.Err
}

func (e *E) Wrap(err error) {
	e.Err = err
	_, e.File, e.Line, _ = runtime.Caller(2)
}

func WrapEf(Err WrapError, format string, a ...interface{}) error {
	var wErr error
	if len(a) > 0 {
		wErr = fmt.Errorf(format, a...)
	} else {
		wErr = errors.New(format)
	}
	Err.Wrap(wErr)
	return Err
}

func WrapE(Err WrapError, err error) error {
	Err.Wrap(Err)
	return Err
}

func NewWrapE(err error) error {
	e := &E{}
	e.Wrap(err)
	return e
}

func NewWrapEf(format string, a ...interface{}) error {
	e := &E{}
	var wErr error
	if len(a) > 0 {
		wErr = fmt.Errorf(format, a...)
	} else {
		wErr = errors.New(format)
	}
	e.Wrap(wErr)
	return e
}

//InternalError
type InternalError struct{ E }
