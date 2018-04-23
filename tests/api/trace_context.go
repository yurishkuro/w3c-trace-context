package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// TraceContext is a container of the trace context fields.
type TraceContext struct {
	TraceParent string
	TraceState  string
}

// ToRequest converts trace context into HTTP headers of the request.
func (c TraceContext) ToRequest(r *http.Request) {
	r.Header.Set("traceparent", c.TraceParent)
	r.Header.Set("tracestate", c.TraceState)
}

// ParseTraceParent breaks the traceparent field into individual components.
func (c TraceContext) ParseTraceParent() (traceID, spanID string, sampled bool, err error) {
	splits := strings.Split(c.TraceParent, "-")
	if len(splits) != 4 {
		err = fmt.Errorf("invalid TraceParent, expecting 4 fields: %s", c.TraceParent)
		return
	}
	if splits[0] != "00" {
		err = fmt.Errorf("invalid TraceParent, expecting version 00: %s", c.TraceParent)
		return
	}
	flags, err1 := strconv.Atoi(splits[3])
	if err1 != nil {
		err = err1
		return
	}
	traceID = splits[1]
	spanID = splits[2]
	sampled = flags != 0
	return
}

// TraceContextFromRequest extracts TraceContext from http.Request.
func TraceContextFromRequest(r *http.Request) TraceContext {
	return TraceContext{
		TraceParent: r.Header.Get("traceparent"),
		TraceState:  r.Header.Get("tracestate"),
	}
}
