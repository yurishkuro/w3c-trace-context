package api

// Span is an abstraction of a vendor span that exposes certain introspection functions.
type Span interface {
	TraceID() string
	SpanID() string
	ParentSpanID() string
	Sampled() bool
	CorrelationID() string

	ToTraceContext() TraceContext
}

// Tracer is an abstraction of a vendor tracer.
type Tracer interface {
	StartSpan(traceContext TraceContext) Span

	Configuration() TracerConfiguration
}
