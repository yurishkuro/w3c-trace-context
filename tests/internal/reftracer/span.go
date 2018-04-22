package reftracer

import (
	"fmt"

	"github.com/w3c/distributed-tracing/tests/api"
)

// Span implements api.Span
type Span struct {
	traceID       string
	spanID        string
	parentSpanID  string
	sampled       bool
	correlationID string
	traceState    string
}

func (s *Span) TraceID() string       { return s.traceID }
func (s *Span) SpanID() string        { return s.spanID }
func (s *Span) ParentSpanID() string  { return s.parentSpanID }
func (s *Span) Sampled() bool         { return s.sampled }
func (s *Span) CorrelationID() string { return s.correlationID }

func (s *Span) ToTraceContext() api.TraceContext {
	sampled := "00"
	if s.sampled {
		sampled = "01"
	}
	return api.TraceContext{
		TraceParent: fmt.Sprintf("00-%s-%s-%s", s.traceID, s.spanID, sampled),
		TraceState:  s.traceState, // TODO encode own vendor position
	}
}
