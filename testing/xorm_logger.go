package testing

import (
	"fmt"
	"testing"
	xorm_core "xorm.io/core"
	"xorm.io/xorm"
)

type XormTestingLogger struct {
	t       *testing.T
	level   xorm_core.LogLevel
	showSQL bool
}

var _ xorm_core.ILogger = (*XormTestingLogger)(nil)

func SetXormTestingLogger(engine *xorm.Engine, t *testing.T) {
	newLogger := NewXormTestingLogger(t)
	oldLogger := engine.Logger()
	newLogger.SetLevel(oldLogger.Level())
	newLogger.ShowSQL(oldLogger.IsShowSQL())
	engine.SetLogger(newLogger)
}

func NewXormTestingLogger(t *testing.T) *XormTestingLogger {
	return &XormTestingLogger{
		t: t,
	}
}

// Error implement core.ILogger
func (l *XormTestingLogger) Error(v ...interface{}) {
	if l.level <= xorm_core.LOG_ERR {
		text := fmt.Sprint(v...)
		l.t.Log("xorm [error]" + text)
	}
	return
}

// Errorf implement core.ILogger
func (l *XormTestingLogger) Errorf(format string, v ...interface{}) {
	if l.level <= xorm_core.LOG_ERR {
		text := fmt.Sprintf(format, v...)
		l.t.Log("xorm [error]" + text)
	}
	return
}

// Debug implement core.ILogger
func (l *XormTestingLogger) Debug(v ...interface{}) {
	if l.level <= xorm_core.LOG_DEBUG {
		text := fmt.Sprint(v...)
		l.t.Log("xorm [debug]" + text)
	}
	return
}

// Debugf implement core.ILogger
func (l *XormTestingLogger) Debugf(format string, v ...interface{}) {
	if l.level <= xorm_core.LOG_DEBUG {
		text := fmt.Sprintf(format, v...)
		l.t.Log("xorm [debug]" + text)
	}
	return
}

// Info implement core.ILogger
func (l *XormTestingLogger) Info(v ...interface{}) {
	if l.level <= xorm_core.LOG_INFO {
		text := fmt.Sprint(v...)
		l.t.Log("xorm [info]" + text)
	}
	return
}

// Infof implement core.ILogger
func (l *XormTestingLogger) Infof(format string, v ...interface{}) {
	if l.level <= xorm_core.LOG_INFO {
		text := fmt.Sprintf(format, v...)
		l.t.Log("xorm [info]" + text)
	}
	return
}

// Warn implement core.ILogger
func (l *XormTestingLogger) Warn(v ...interface{}) {
	if l.level <= xorm_core.LOG_WARNING {
		text := fmt.Sprint(v...)
		l.t.Log("xorm [warn]" + text)
	}
	return
}

// Warnf implement core.ILogger
func (l *XormTestingLogger) Warnf(format string, v ...interface{}) {
	if l.level <= xorm_core.LOG_WARNING {
		text := fmt.Sprintf(format, v...)
		l.t.Log("xorm [warn]" + text)
	}
	return
}

// Level implement core.ILogger
func (l *XormTestingLogger) Level() xorm_core.LogLevel {
	return l.level
}

// SetLevel implement core.ILogger
func (l *XormTestingLogger) SetLevel(level xorm_core.LogLevel) {
	l.level = level
	return
}

// ShowSQL implement core.ILogger
func (l *XormTestingLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (l *XormTestingLogger) IsShowSQL() bool {
	return l.showSQL
}
