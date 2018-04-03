package tracer

import "github.com/w3c/distributed-tracing/tests/api"

var _ api.Tracer = new(Tracer)
var _ api.Span = new(Span)
