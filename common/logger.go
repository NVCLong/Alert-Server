package common

import (
	"fmt"
	"time"
)

type AbstractLogger interface {
	Debug(message string)
	Verbose(message string)
	Error(message string)
	Log(message string)
	GetTraceId() any
}

type TracingLogger struct {
	context   string
	tracingId any
}

func (s *TracingLogger) GetTraceId() any {
	return s.tracingId
}

func NewTracingLogger(context string, tracingId any) AbstractLogger {
	return &TracingLogger{
		context:   context,
		tracingId: tracingId,
	}
}

func (s *TracingLogger) Debug(message string) {
	time := time.Now().Format("2006/01/02 - 15:04:05")
	fmt.Printf("[GIN-Debug] %s |%s| |%s| - %s\n", time, s.tracingId, s.context, message)
}

func (s *TracingLogger) Log(message string) {
	time := time.Now().Format("2006/01/02 - 15:04:05")
	fmt.Printf("[GIN-Log] %s |%s| |%s| - %s\n", time, s.tracingId, s.context, message)

}

func (s *TracingLogger) Verbose(message string) {
	time := time.Now().Format("2006/01/02 - 15:04:05")
	fmt.Printf("[GIN-Verbose] %s |%s| |%s| - %s\n", time, s.tracingId, s.context, message)
}

func (s *TracingLogger) Error(message string) {
	time := time.Now().Format("2006/01/02 - 15:04:05")
	fmt.Printf("[GIN-Error] %s |%s| |%s| - %s\n", time, s.tracingId, s.context, message)
}
