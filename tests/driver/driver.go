package driver

import (
	"log"
	"net/http"

	"github.com/crossdock/crossdock-go"

	"github.com/w3c/distributed-tracing/tests/driver/behaviors/diffvendor"
	"github.com/w3c/distributed-tracing/tests/driver/behaviors/malformed"
	"github.com/w3c/distributed-tracing/tests/driver/behaviors/missing"
	"github.com/w3c/distributed-tracing/tests/driver/behaviors/samevendor"
	"github.com/w3c/distributed-tracing/tests/driver/params"
)

var behaviors = crossdock.Behaviors{
	params.BehaviorMalformedTraceContext:  malformed.Execute,
	params.BehaviorMissingTraceContext:    missing.Execute,
	params.BehaviorTraceContextSameVendor: samevendor.Execute,
	params.BehaviorTraceContextDiffVendor: diffvendor.Execute,
}

// Start registers behaviors and begins the Crossdock test driver.
func Start() {
	http.Handle("/", crossdock.Handler(behaviors, true))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
