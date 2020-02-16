package log

import (
	"os"
	"sync"
)

var (
	std = New(nil, nil)
)

// Logger Logger struct
type Logger struct {
	mu       sync.Mutex
	buf      []byte
	output   OutRouter
	formater Formater
}

//Output out put log
func (l *Logger) Output(data *Log) (err error) {
	output := l.output.To(data)
	if output == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buf = l.buf[:0]
	l.formater.Format(&l.buf, data)

	_, err = output.Write(l.buf)
	return err
}

//SetOutRouter set the OutRouter of the logger
func (l *Logger) SetOutRouter(newout OutRouter) {
	l.output = newout
}

//New create a new logger
func New(out OutRouter, formater Formater) *Logger {
	if formater == nil {
		formater = DefaultFormater{}
	}
	if out == nil {
		out = NewDefaultOutRouter(os.Stdout, os.Stderr)
	}
	return &Logger{output: out, formater: formater}
}
