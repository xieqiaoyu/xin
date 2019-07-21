package log

import (
	"testing"
)

func BenchmarkWriteErrorWithTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write(Error, true, "This is an err :%s", "just for benchemark")
	}
}

func BenchmarkWriteErrorWithNoTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write(Error, false, "This is an err :%s", "just for benchemark")
	}
}
