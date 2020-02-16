package log

import (
	"time"
)

//Formater formater tells how to format Log into byte buffer
type Formater interface {
	Format(buf *[]byte, l *Log)
}

//DefaultFormater the default formater xin log use
type DefaultFormater struct{}

func (DefaultFormater) Format(buf *[]byte, data *Log) {
	*buf = append(*buf, time.Now().Format(time.RFC3339)...)
	*buf = append(*buf, ' ')
	*buf = append(*buf, '[')
	*buf = append(*buf, data.Level.String()...)
	*buf = append(*buf, ']')
	if data.Tag != "" {
		*buf = append(*buf, ' ')
		*buf = append(*buf, '[')
		*buf = append(*buf, data.Tag...)
		*buf = append(*buf, ']')
	}
	if len(data.Msg) > 0 {
		*buf = append(*buf, ' ')
		*buf = append(*buf, data.Msg...)
	}
	if data.Caller != nil {
		*buf = append(*buf, ' ')
		*buf = append(*buf, '[')
		*buf = append(*buf, data.Caller.File...)
		*buf = append(*buf, ':')
		itoa(buf, data.Caller.Line, -1)
		*buf = append(*buf, ']')
	}
	*buf = append(*buf, '\n')
}

//itoa copy from go log package
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
