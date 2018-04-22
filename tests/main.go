// Copyright

package main

import (
	"github.com/w3c/distributed-tracing/tests/actor"
	"github.com/w3c/distributed-tracing/tests/driver"
	"github.com/w3c/distributed-tracing/tests/internal/reftracer"
)

func main() {
	actor.New(reftracer.New()).Start()
	driver.Start()
}
