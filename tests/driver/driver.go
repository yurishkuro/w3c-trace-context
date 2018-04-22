package driver

import (
	"log"
	"net/http"

	"github.com/crossdock/crossdock-go"

	"github.com/w3c/distributed-tracing/tests/driver/behaviors/trace"
	"github.com/w3c/distributed-tracing/tests/driver/params"
)

var behaviors = crossdock.Behaviors{
	params.BehaviorMalformedTraceContext:  trace.Trace,
	params.BehaviorNoTraceContext:         trace.Trace,
	params.BehaviorTraceContextSameVendor: trace.Trace,
	params.BehaviorTraceContextDiffVendor: trace.Trace,
}

// Start registers behaviors and begins the Crossdock test driver.
func Start() {
	http.Handle("/", crossdock.Handler(behaviors, true))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
