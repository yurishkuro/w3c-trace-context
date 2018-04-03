package driver

import (
	"log"
	"net/http"

	"github.com/crossdock/crossdock-go"

	"github.com/w3c/distributed-tracing/tests/driver/behaviors/trace"
)

var behaviors = crossdock.Behaviors{
	"trace": trace.Trace,
}

// Start registers behaviors and begins the Crossdock test driver.
func Start() {
	http.Handle("/", crossdock.Handler(behaviors, true))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
