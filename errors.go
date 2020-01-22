package xin

import (
	"errors"
	"fmt"
	"runtime"
)

//WrapError WrapError
type WrapError interface {
	error
	//Wrap Wrap a error in current error
	Wrap(error)
}

//TracedError a WrapError store the position where error occur
type TracedError struct {
	File string
	Line int
	Err  error
}

//Error error interface
func (e *TracedError) Error() string {
	return fmt.Sprintf("%s (%s:%d)", e.Err, e.File, e.Line)
}

//Unwrap use for errors.Unwrap
func (e *TracedError) Unwrap() error {
	return e.Err
}

//Wrap WrapError interface,do not call this function directly , use  WrapE func instead
func (e *TracedError) Wrap(err error) {
	e.Err = err
	_, e.File, e.Line, _ = runtime.Caller(2)
}

//WrapE wrap the given into the given WrapError
func WrapE(Err WrapError, err error) error {
	Err.Wrap(Err)
	return Err
}

//WrapEf create an error by format string and wrap it into the given WrapError
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

//NewWrapE @deprecated use NewTracedE instead
// create a new TracedError and wrap the given error into it
func NewWrapE(err error) error {
	e := &TracedError{}
	e.Wrap(err)
	return e
}

//NewTracedE create a new TracedError and wrap the given error into it
func NewTracedE(err error) error {
	e := &TracedError{}
	e.Wrap(err)
	return e
}

//NewWrapEf @deprecated use NewTracedEf instead
// create an error by format string and wrap it into a new TracedError
func NewWrapEf(format string, a ...interface{}) error {
	e := &TracedError{}
	var wErr error
	if len(a) > 0 {
		wErr = fmt.Errorf(format, a...)
	} else {
		wErr = errors.New(format)
	}
	e.Wrap(wErr)
	return e
}

//NewTracedEf create an error by format string and wrap it into a new TracedError
func NewTracedEf(format string, a ...interface{}) error {
	e := &TracedError{}
	var wErr error
	if len(a) > 0 {
		wErr = fmt.Errorf(format, a...)
	} else {
		wErr = errors.New(format)
	}
	e.Wrap(wErr)
	return e
}

//InternalError an error for framework internal use
type InternalError struct{ TracedError }
