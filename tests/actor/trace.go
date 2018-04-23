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

// Trace implements actor's "trace" endpoint. The main responsibility of this endpoint are:
//   - read trace context fields from the headers
//   - ask the tracer implementation to parse them into a "span"
//   - if the request contains an instruction to call downstream node, then
//     - ask tracer to encode the new span into trace context fields
//     - execute the request against the downstream node
//     - attach the response to the overall response
//   - record the attributes of the tracer's "span" in the overall response
//   - record the tracer config in the overall response
//   - return the response
func (a *Actor) Trace(w http.ResponseWriter, r *http.Request) {
	var req api.Request
	xhttp.HandleJSON(w, r, &req, func(r *http.Request, in interface{}) (interface{}, error) {
		{ // log request
			jsonBytes, _ := json.Marshal(in)
			log.Printf("actor %s received request: %s", a.name, string(jsonBytes))
		}
		if a.name != req.Actor {
			return nil, fmt.Errorf("Current actor name '%s' does not match target actor name '%s'", a.name, req.Actor)
		}

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
				TraceID:       span.TraceID(),
				SpanID:        span.SpanID(),
				Sampled:       span.Sampled(),
				ParentSpanID:  span.ParentSpanID(),
				CorrelationID: span.CorrelationID(),
				TraceParent:   tc.TraceParent,
				TraceState:    tc.TraceState,
			},
			TracerConfig: a.tracer.Configuration(),
			Downstream:   downstream,
		}, nil
	})
}

func (a *Actor) callDownstream(dn *api.Request, span api.Span) (*api.Response, error) {
	if dn.Actor == "" {
		return nil, fmt.Errorf("no actor name")
	}
	server := dn.Actor
	if dn.Server != "" {
		server = dn.Server
	}
	url := fmt.Sprintf("http://%s:8081/trace", server)
	log.Printf("calling downstream %s at %s", dn.Actor, url)
	tc := span.ToTraceContext()
	var res api.Response
	err := xhttp.PostJSON(
		context.Background(),
		url,
		tc.ToRequest,
		&api.Request{
			Actor:      dn.Actor,      // for sanity checks
			Server:     server,        // for unit tests
			Downstream: dn.Downstream, // pass recursive downstream call instructions, if any
		},
		&res,
	)
	return &res, err
}
