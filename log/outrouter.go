package log

import (
	"io"
)

//OutRouter outrouter decide where to log
type OutRouter interface {
	To(l *Log) io.Writer
}

//DefaultOutRouter DefaultOutRouter
type DefaultOutRouter struct {
	out    io.Writer
	errout io.Writer
}

//To Implement OutRouter
func (r *DefaultOutRouter) To(l *Log) io.Writer {
	if l.Level.GetCode() < Warning.GetCode() {
		return r.out
	}
	return r.errout
}

//NewDefaultOutRouter NewDefaultOutRouter
func NewDefaultOutRouter(out, errout io.Writer) *DefaultOutRouter {
	return &DefaultOutRouter{
		out:    out,
		errout: errout,
	}
}
