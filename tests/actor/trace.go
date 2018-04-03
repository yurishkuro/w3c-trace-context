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

// Trace implements actor's "trace" behavior.
func (a *Actor) Trace(w http.ResponseWriter, r *http.Request) {
	var req api.Request
	xhttp.HandleJSON(w, r, &req, func(r *http.Request, in interface{}) (interface{}, error) {
		jsonBytes, _ := json.Marshal(in)
		log.Printf("actor received request: %s", string(jsonBytes))

		tc := api.TraceContextFromRequest(r)
		log.Printf("received trce context: %+v", tc)
		span := a.tracer.StartSpan(tc)

		var downstream *api.Response
		if req.Downstream != nil {
			d, err := a.CallDownstream(req.Downstream, span)
			if err != nil {
				return nil, err
			}
			downstream = d

			jsonBytes, _ := json.Marshal(d)
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

func (a *Actor) CallDownstream(dn *api.DownstreamRequest, span api.Span) (*api.Response, error) {
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
