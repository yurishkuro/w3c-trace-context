# Compatibility Test-bed

This module contains a test harness that can be used to verify a given tracer's compliance and compatibility with the spec.

## Getting started

* Have a Go toolchain installed
* `go get github.com/crossdock/crossdock-go` (a dependency)
* Clone this repo to `$GOPATH/src/github.com/w3c/distributed-tracing/`;
  * `mkdir -p $GOPATH/src/github.com/w3c/`
  * `cd $GOPATH/src/github.com/w3c/`
  * `git clone git@github.com:yurishkuro/distributed-tracing.git`
  * `cd distributed-tracing`
  * `git checkout compliance-tests`
* Run unit tests: `make test`
* Run actual test bed: `make crossdock`

## Test-bed Components

### Orchestrator

The Crossdock framework is used as orchestrator. At the moment it's implemented only in the unit tests `main_test.go`,
but can be easily configured to run via docker-compose.

### Driver

Driver is a binary that receives request from orchestrator with instructions about a specific test to perform. A sample request looks like this:

```
GET http://127.0.0.1:8080?actor1=ref&actor2=ref&behavior=trace&sampled=true
```

where:
  * `behavior` is the name of the test to be executed with given parameters (see Behaviors below)
  * `actor1` is the name of the first Node in the chain of RPC calls (`ref` means use reference implementation)
  * `actor1` is the name of the second Node in the chain of RPC calls (`ref` means use reference implementation)
  * `sampled` tells the driver to initiate a trace in sampled state

### Actors

Actors implement Nodes in the test case that exchange RPC requests. To avoid having each vendor re-implement the exact behavior expected of the Node the default implementation is generic enough so that any vendor tracer can be plugged in by implementing the `api.Tracer` interface.

## Behaviors

The driver and actors can support different behaviors (types of test). Currently a single behavior `trace` is implemented.

### Behavior "Trace"

#### Parameters

* actor1 - name of the first node called by the driver
* actor2 - name of the second node called by the first node
* sampled - whether the driver sends a trace context as sampled

#### Driver

When executing this test, the driver

* manufactures a new trace and encodes it in Trace-Parent
* populates Trace-State with fake vendor entries
* creates a request to actor1 with instructions to call the second actor (by providing its name)
* upon receiving the response from actor1 validates that both actors observed expected trace context headers with expected causal relationships between spans

#### Actors

* Actors implement `/trace` endpoint
* They return a response that records the trace/span IDs of the span for that node, and other fields
* If request contains an instruction to call another actor, the first actor executes it and embeds the other actor's response into its own response

## TODO

### Dimensions of the individual tests
  * permutations of the inbound trace context
    * malformed trace context (many variations)
    * sampled or unsampled
    * different versions (in the future)
    * inbound trace context by different vendor
    * inbound trace context by the same vendor (requires 2 actor hops + validation)
  * participation mode of the Node
    * trust
    * no trust - different vendor
    * no trust - same vendor
  * sampling - how does the Node decide on sampling given
    * no inbound trace context
    * malformed trace context
    * unsampled trace context
      * Node keeps no sampling
      * Node up-samples

Node parameters:
  * TRUST_TRACE_ID = true/false
  * TRUST_SAMPLING = true/false
  * SAMPLE = true/false (when no inbound trace context)
  * UPSAMPLE = true/false (when inbound trace context is not sampled)

How does the driver know about Node's parameters?
  * Make the Node return them in the response
  * Pass parameters to the driver (requires too many compose files)


* Use a real tracer that supports Trace Context semantics to implement the `api.Trace` and create another actor
* Define docker-compose file that can pull different types of actors into one test suite
* Decide how different participation modes should be modelled
  * Option 1 - as different behaviors. Downside of that is it requires a separate endpoint in each actor, although most of the code is similar
  * Option 2 - as additional parameters of the `trace` behavior. For example:
    * `actor1_participation = pass-through | join | correlate`
    * `actor1_participation = pass-through | join | correlate`
