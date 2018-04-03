// Copyright

package main

import (
	"github.com/w3c/distributed-tracing/tests/actor"
	"github.com/w3c/distributed-tracing/tests/driver"
)

func main() {
	actor.New(nil).Start()
	driver.Start()
}
