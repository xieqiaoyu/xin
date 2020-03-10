package log

import (
	"bytes"
	"log"
	"testing"
)

func BenchmarkGoLogWithTrace(b *testing.B) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	for i := 0; i < b.N; i++ {
		log.Printf("[%s] This is an err :%s", "Error", "just for benchemark")
	}
}
func BenchmarkWriteErrorWithTrace(b *testing.B) {
	buf := new(bytes.Buffer)
	Std.SetOutRouter(NewDefaultOutRouter(buf, buf))
	for i := 0; i < b.N; i++ {
		Write(Error, 1, true, "This is an err :%s", "just for benchemark")
	}
}

func BenchmarkGoLogWithNoTrace(b *testing.B) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	log.SetFlags(log.LstdFlags)
	for i := 0; i < b.N; i++ {
		log.Printf("[%s] This is an err :%s", "Error", "just for benchemark")
	}
}
func BenchmarkWriteErrorWithNoTrace(b *testing.B) {
	buf := new(bytes.Buffer)
	Std.SetOutRouter(NewDefaultOutRouter(buf, buf))
	for i := 0; i < b.N; i++ {
		Write(Error, 1, false, "This is an err :%s", "just for benchemark")
	}
}

func TestLogOutput(t *testing.T) {
	Debugf("debug")
	Infof("info")
	Noticef("notice")
	Warningf("warning")
	Errorf("error")
	Alertf("alert")
	Criticalf("critical")
	Emergencyf("3mergency")
}
