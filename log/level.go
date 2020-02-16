package log

import (
	"fmt"
)

type LogLevel interface {
	GetCode() int
	fmt.Stringer
}

var levelLiteral = map[RFC5424SeverityLevel]string{
	Debug:     "Debug",
	Info:      "Info",
	Notice:    "Notice",
	Warning:   "Warning",
	Error:     "Error",
	Critical:  "Critical",
	Alert:     "Alert",
	Emergency: "Emergency",
}

//RFC5424SeverityLevel  RFC5424SeverityLevel
type RFC5424SeverityLevel int

//GetCode logLevel interface
func (s RFC5424SeverityLevel) GetCode() int {
	return int(s)
}

//GetLiteral logLevel interface
func (s RFC5424SeverityLevel) String() string {
	if literal, ok := levelLiteral[s]; ok {
		return literal
	}
	return "Unknown"
}

var _ LogLevel = RFC5424SeverityLevel(0)

// predefine log level
const (
	Debug     RFC5424SeverityLevel = 7
	Info      RFC5424SeverityLevel = 6
	Notice    RFC5424SeverityLevel = 5
	Warning   RFC5424SeverityLevel = 4
	Error     RFC5424SeverityLevel = 3
	Critical  RFC5424SeverityLevel = 2
	Alert     RFC5424SeverityLevel = 1
	Emergency RFC5424SeverityLevel = 0
)
