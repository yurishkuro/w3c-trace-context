// Copyright

package main

import (
	"github.com/w3c/distributed-tracing/tests/actor"
	"github.com/w3c/distributed-tracing/tests/driver"
	"github.com/w3c/distributed-tracing/tests/internal/reftracer"
)

// This is a sample main for implementing an actor using vendor-specific tracer.
func main() {
	// The real vendor code will implement their own api.Tracer
	tracer := reftracer.New()
	actor.New(tracer).Start()
	driver.Start()
}
