package actor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/w3c/distributed-tracing/tests/api"
	"github.com/w3c/distributed-tracing/tests/internal/xhttp"
)

// Trace implements actor's "trace" endpoint.
func (a *Actor) Trace(w http.ResponseWriter, r *http.Request) {
	var req api.Request
	xhttp.HandleJSON(w, r, &req, func(r *http.Request, in interface{}) (interface{}, error) {
		jsonBytes, err := json.Marshal(in)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse request as JSON: %v", err)
		}
		log.Printf("actor received request: %s", string(jsonBytes))

		tc := api.TraceContextFromRequest(r)
		log.Printf("received trace context: %+v", tc)
		span := a.tracer.StartSpan(tc)

		var downstream *api.Response
		if req.Downstream != nil {
			d, err := a.callDownstream(req.Downstream, span)
			if err != nil {
				return nil, err
			}
			downstream = d

			jsonBytes, err := json.Marshal(d)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse downstream response as JSON: %v", err)
			}
			log.Printf("downstream response: %s", string(jsonBytes))
		}

		return &api.Response{
			Trace: api.ObservedTrace{
				TraceID:      span.TraceID(),
				SpanID:       span.SpanID(),
				Sampled:      span.Sampled(),
				ParentSpanID: span.ParentSpanID(),
				TraceParent:  tc.TraceParent,
				TraceState:   tc.TraceState,
			},
			Downstream: downstream,
		}, nil
	})
}

func (a *Actor) callDownstream(dn *api.Request, span api.Span) (*api.Response, error) {
	if dn.Actor == "" {
		return nil, fmt.Errorf("no actor name")
	}
	log.Printf("calling downstream %s", dn.Actor)
	tc := span.ToTraceContext()
	var res api.Response
	err := xhttp.PostJSON(
		context.Background(),
		"http://127.0.0.1:8081/trace", // TODO use target actor host name
		tc.ToRequest,
		&api.Request{
			Downstream: dn.Downstream,
		},
		&res,
	)
	return &res, err
}
