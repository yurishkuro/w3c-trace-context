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
	CorrelationID string `json:"correlation_id,omitempty"` // TODO should this contain the full traceparent header?
	TraceParent   string `json:"trace_parent,omitempty"`
	TraceState    string `json:"trace_state,omitempty"`
}
