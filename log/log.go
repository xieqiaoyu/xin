package log

import (
	"fmt"
	"runtime"
)

//WithTag new Log with tag
func WithTag(tag string) *Log {
	return &Log{
		Tag: tag,
	}
}

//Log one line log
type Log struct {
	Level  LogLevel
	Tag    string
	Msg    string
	Caller *Caller
}

//Caller where the log is call from
type Caller struct {
	File string
	Line int
}

//Debugf Debugf
func Debugf(format string, v ...interface{}) {
	Write(Debug, 2, true, format, v...)
}

//Debugf Debugf
func (l *Log) Debugf(format string, v ...interface{}) {
	l.Write(Debug, 2, true, format, v...)
}

//Infof Infof
func Infof(format string, v ...interface{}) {
	Write(Info, 2, false, format, v...)
}

//Infof Infof
func (l *Log) Infof(format string, v ...interface{}) {
	l.Write(Info, 2, false, format, v...)
}

//Noticef Noticef
func Noticef(format string, v ...interface{}) {
	Write(Notice, 2, false, format, v...)
}

//Noticef Noticef
func (l *Log) Noticef(format string, v ...interface{}) {
	l.Write(Notice, 2, false, format, v...)
}

//Warningf Warningf
func Warningf(format string, v ...interface{}) {
	Write(Warning, 2, true, format, v...)
}

//Warningf Warningf
func (l *Log) Warningf(format string, v ...interface{}) {
	l.Write(Warning, 2, true, format, v...)
}

//Errorf Errorf
func Errorf(format string, v ...interface{}) {
	Write(Error, 2, true, format, v...)
}

//Errorf Errorf
func (l *Log) Errorf(format string, v ...interface{}) {
	l.Write(Error, 2, true, format, v...)
}

//Criticalf Criticalf
func Criticalf(format string, v ...interface{}) {
	Write(Critical, 2, true, format, v...)
}

//Criticalf Criticalf
func (l *Log) Criticalf(format string, v ...interface{}) {
	l.Write(Critical, 2, true, format, v...)
}

//Alertf Alertf
func Alertf(format string, v ...interface{}) {
	Write(Alert, 2, true, format, v...)
}

//Alertf Alertf
func (l *Log) Alertf(format string, v ...interface{}) {
	l.Write(Alert, 2, true, format, v...)
}

//Emergencyf Emergencyf
func Emergencyf(format string, v ...interface{}) {
	Write(Emergency, 2, true, format, v...)
}

//Emergencyf Emergencyf
func (l *Log) Emergencyf(format string, v ...interface{}) {
	l.Write(Emergency, 2, true, format, v...)
}

//Write write log
func (l Log) Write(level LogLevel, calldepth int, trace bool, format string, v ...interface{}) {
	l.Level = level
	l.Msg = fmt.Sprintf(format, v...)
	if trace {
		_, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.Caller = &Caller{
			File: file,
			Line: line,
		}
	}
	std.Output(&l)
}

//Write write log
func Write(level LogLevel, calldepth int, trace bool, format string, v ...interface{}) {
	l := &Log{
		Level: level,
		Msg:   fmt.Sprintf(format, v...),
	}
	if trace {
		_, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.Caller = &Caller{
			File: file,
			Line: line,
		}
	}
	std.Output(l)
}
