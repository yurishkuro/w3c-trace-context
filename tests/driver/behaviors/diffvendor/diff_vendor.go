package diffvendor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
	log.Printf("executing behavior %s", params.BehaviorTraceContextDiffVendor)
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
				Actor:  params.RefActor,
				Server: bp.server,
			},
		},
		&res,
	)
	if fatals.NoError(err, "failed to post JSON") {
		json, _ := json.Marshal(&res)
		log.Printf("driver received response: %s", json)
	}

	assert := crossdock.Assert(t)
	if res.TracerConfig.TrustTraceID {
		assert.Equal(traceID, res.Trace.TraceID, "same trace ID")
	} else {
		assert.NotEqual(traceID, res.Trace.TraceID, "different trace ID")
		assert.Equal(traceID, res.Trace.CorrelationID, "trace ID is in correlationID")
	}
	assert.NotEmpty(res.Trace.SpanID, "spanID is not empty")
	assert.Equal(spanID, res.Trace.ParentSpanID, "ParentSpanID equal root spanID")
	assert.Equal(true, res.Trace.Sampled, "span is sampled")

	// downstream validation
	fatals.NotNil(res.Downstream, "downstream response not empty")
	if res.TracerConfig.TrustTraceID {
		assert.Equal(traceID, res.Downstream.Trace.TraceID, "same downstream traceID")
	} else {
		assert.Equal(res.Trace.TraceID, res.Downstream.Trace.TraceID, "downstream traceID equal 1st actor's traceID")
	}
	assert.Equal(true, res.Downstream.Trace.Sampled, "downstream span is sampled")
	// validate vendor key in the 1st position
	assert.NotEqual(tc.TraceState, res.Downstream.Trace.TraceState, "modified tracestate")
	assert.NotEmpty(res.TracerConfig.VendorKey, "non-empty vendor key")
	vendorParts := strings.Split(res.Downstream.Trace.TraceState, ",")
	firstParts := strings.Split(vendorParts[0], "=")
	assert.Equal(res.TracerConfig.VendorKey, firstParts[0], "vendor key '%s' in the first position", res.TracerConfig.VendorKey)
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
