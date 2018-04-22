package api

// Request is sent by the driver to the first Node, and by Nodes to downstream Nodes.
type Request struct {
	Actor      string   `json:"actor"`
	Server     string   `json:"server,omitempty"` // only used in unit tests to override IP address of the actor
	Downstream *Request `json:"downstream,omitempty"`
}

// Response is the response of the Node.
type Response struct {
	TracerConfig TracerConfiguration `json:"tracer_config"`
	Trace        ObservedTrace       `json:"trace"`
	Downstream   *Response           `json:"downstream,omitempty"`
}

// ObservedTrace describes the trace that the node observed / recorded.
type ObservedTrace struct {
	TraceID       string `json:"trace_id,omitempty"`
	SpanID        string `json:"span_id,omitempty"`
	ParentSpanID  string `json:"parent_id,omitempty"`
	Sampled       bool   `json:"sampled,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	TraceParent   string `json:"trace_parent,omitempty"`
	TraceState    string `json:"trace_state,omitempty"`
}

// TracerConfiguration describes how the actor's tracer is going to behave under different conditions.
type TracerConfiguration struct {
	ActorName string

	// TrustTraceID controls whether the tracer respects inbound trace ID or creates a new trace
	// and records inbound trace ID as correlation.
	TrustTraceID bool

	// TrustSampling control whether the tracer respects inbound sampling flag or makes its own decision (based on Sample below).
	TrustSampling bool

	// Sample controls which sampling decision the tracer makes when it needs to make it (e.g when there is no inbound trace context).
	Sample bool

	// Upsample controls whether the tracer will switch on sampling even if the inbound trace context has sampling=off.
	Upsample bool
}
