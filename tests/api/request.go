package api

type DownstreamRequest struct {
	Actor      string             `json:"actor"`
	Downstream *DownstreamRequest `json:"downstream,omitempty"`
}

type Request struct {
	Downstream *DownstreamRequest `json:"downstream,omitempty"`
}

type ObservedTrace struct {
	TraceID       string `json:"trace_id,omitempty"`
	SpanID        string `json:"span_id,omitempty"`
	ParentSpanID  string `json:"parent_id,omitempty"`
	Sampled       bool   `json:"sampled,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	TraceParent   string `json:"trace_parent,omitempty"`
	TraceState    string `json:"trace_state,omitempty"`
}

// type DownstreamResponse struct {
// 	Actor      string              `json:"actor"`
// 	Trace      ObservedTrace       `json:"trace"`
// 	Downstream *DownstreamResponse `json:"downstream,omitempty"`
// }

type Response struct {
	Actor      string        `json:"actor"`
	Trace      ObservedTrace `json:"trace"`
	Downstream *Response     `json:"downstream,omitempty"`
}
