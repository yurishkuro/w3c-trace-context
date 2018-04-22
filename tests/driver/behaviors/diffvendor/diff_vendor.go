package diffvendor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	crossdock "github.com/crossdock/crossdock-go"
	"github.com/w3c/distributed-tracing/tests/api"
	"github.com/w3c/distributed-tracing/tests/driver/params"
	"github.com/w3c/distributed-tracing/tests/internal/random"
	"github.com/w3c/distributed-tracing/tests/internal/xhttp"
)

// TraceParentVersion is the version of the spec.
const TraceParentVersion = "00"

type behaviorParams struct {
	actor  string
	server string
}

// Execute implements the 'trace-context-diff-vendor' behavior.
func Execute(t crossdock.T) {
	fatals := crossdock.Fatals(t)
	bp := readParams(t)
	log.Printf("params %+v", bp)

	traceID := random.New64BitID() + random.New64BitID()
	spanID := random.New64BitID()
	flags := "01"
	tc := api.TraceContext{
		TraceParent: fmt.Sprintf("%s-%s-%s-%s", TraceParentVersion, traceID, spanID, flags),
		TraceState:  "vnd1=abcd,vnd2=xyz",
	}

	server := bp.actor
	if bp.server != "" {
		server = bp.server
	}
	url := fmt.Sprintf("http://%s:8081/trace", server)

	var res api.Response
	err := xhttp.PostJSON(
		context.Background(),
		url,
		tc.ToRequest,
		&api.Request{
			Actor: bp.actor,
			Downstream: &api.Request{
				Actor: params.RefActor,
			},
		},
		&res,
	)
	if fatals.NoError(err, "failed to post JSON") {
		json, _ := json.Marshal(&res)
		log.Printf("driver received response: %s", json)
	}

	assert := crossdock.Assert(t)
	assert.Equal(traceID, res.Trace.TraceID, "same trace ID")
	assert.NotEmpty(res.Trace.SpanID)
	assert.Equal(spanID, res.Trace.ParentSpanID)
	assert.Equal(true, res.Trace.Sampled)
	assert.Equal(tc.TraceParent, res.Trace.TraceParent)
	assert.Equal(tc.TraceState, res.Trace.TraceState)

	fatals.NotNil(res.Downstream)
	assert.Equal(traceID, res.Downstream.Trace.TraceID)
	assert.Equal(res.Trace.SpanID, res.Downstream.Trace.ParentSpanID)
	assert.Equal(true, res.Downstream.Trace.Sampled)
	assert.Equal(tc.TraceState, res.Downstream.Trace.TraceState)
}

func readParams(t crossdock.T) behaviorParams {
	fatals := crossdock.Fatals(t)

	b := behaviorParams{
		actor:  t.Param(params.Actor),
		server: t.Param(params.Server),
	}
	fatals.NotEmpty(b.actor, "actor cannot be empty")
	return b
}
