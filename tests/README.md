# Compatibility Test Suite

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
* Run actual test suite: `make crossdock`

## Test Suite Components

### Orchestrator

The Crossdock framework is used as orchestrator. It is invoked from the [docker-compose](./docker-compose.yaml) file.

### Driver

Driver is a binary that receives request from orchestrator with instructions about a specific test to perform. A sample request looks like this:

```
GET http://127.0.0.1:8080?actor=refnode&behavior=trace_context_diff_vendor
```

where:
  * `behavior` is the name of the test to be executed with given parameters (see Behaviors below)
  * `actor` is the name of the Node being tested (`refnode` service provide reference implementation)

The driver is implemented by the `refnode` service (see [docker-compose.yaml](./docker-compose.yaml)).

### Actors

Actors implement Nodes in the test case that exchange RPC requests. To avoid having each vendor re-implement the exact behavior expected of the Node the default implementation is generic enough so that any vendor tracer can be plugged in by implementing the `api.Tracer` interface.

The reference implementation of the actor is implemented by the `refnode` service (see [docker-compose.yaml](./docker-compose.yaml)). It can be configured to run in different participation modes, e.g. the same code is used for service `refnode1` which is configured to not trust the inbound trace ID and always restart the trace.

So far the [Actor module](./actor/)  only implements a single endpoint `/trace`. See Actor struct comments and the [request/response API](./api/). Actors return a response that records the trace/span IDs of the span for that node, and other fields. If request contains an instruction to call another actor, the first actor executes it and embeds the other actor's response into its own response

## Behaviors

The driver supports different behaviors (types of test) defined in [driver/behaviors](./driver/behaviors/) package.

### Behavior "malformed_trace_context"

Tests how the actor reacts to a malformed trace context headers. Currently not implemented.

RPC chain: `driver->vendor->refnode`.

### Behavior "missing_trace_context"

Tests how the actor reacts to missing trace context headers. Currently not implemented.

RPC chain: `driver->vendor->refnode`.

### Behavior "trace_context_diff_vendor"

Tests how the actor reacts to well-formed trace context by different vendors.

RPC chain: `driver->vendor->refnode`.

When executing this test, the driver

* manufactures a new trace and encodes it in Trace-Parent;
* populates Trace-State with fake vendor entries;
* creates a request to the `actor` service with instructions to call the second `refnode` actor, which can record outbound trace context;
* upon receiving the response from the main actor validates that both actors observed expected trace context headers with expected causal relationships between spans.

### Behavior "trace_context_diff_vendor"

Tests how the actor reacts to well-formed trace context from the same vendor. Currently not implemented.

RPC chain: `driver->vendor->vendor->refnode` (because the driver would not know how to prepare the first trace context with the correct vendor key).

## How to test vendor-specific implementation for Trace-Context compatibility.

The [docker-compose.yaml](./docker-compose.yaml) file uses `example1` container as a substitute for a vendor-provided container. The basic steps for a vendor are the following (Go only, for now):
  * implement `api.Tracer` interface
  * create a binary similar to [example/main.go](./example/main.go) that uses default Actor implementation with its own Tracer
  * create a Docker image from the binary (see [example/Dockerfile](./example/Dockerfile))
  * update the main [docker-compose.yaml](./docker-compose.yaml) file to run the new image as a service, similar to `example1`
    * the container can be used multiple times with different environment variables, similar to `refnode` and `refnode1`

To test implementations in other languages, this test suite needs to implement a reusable Actor, so that vendors would only need to provide the `Tracer` implementation.

Another alternative is to implement Actor completely independently, perhaps using vendor's own instrumentation, using the spec below.

## Actor Specification

The exact data types are defined in [api/actor.go](./api/actor.go).

### Incoming Requests

The first actor in the chain usually receives `HTTP POST http://actor-name:8081/trace` with JSON payload that looks like this:

```
{
  "actor": "actor-name",
  "downstream": {
    "actor": "next-actor-name"
  }
}
```

### Downstream Requests

The `downstream` portion is optional, if present the actor is expected to call the next actor and pass it just that downstream portion, i.e. in this case `HTTP POST http://next-actor-name:8081/trace` with JSON payload 

```
{
  "actor":"next-actor-name"
}
```

### Response

The response contains three parts:
  1. the description of the server-side span created by the actor
  1. the description of actor's configuration parameters that explain its behavior
  1. if request had the `downstream` section, the response must include the response from the downstream actor.

Below is an example of the top-level actor's response:

```
{
  "tracer_config": {
    "ActorName": "actor-name",
    "VendorKey": "someKey",
    "TrustTraceID": true,
    "TrustSampling": true,
    "Sample": true,
    "Upsample": false
  },
  "trace": {
    "trace_id": "999ed0ff376053d0ce00566d456e9c80",
    "span_id": "87f1df57ec17c593",
    "parent_id": "005ba7afcb2da550",
    "sampled": true,
    "trace_parent": "00-999ed0ff376053d0ce00566d456e9c80-005ba7afcb2da550-01",
    "trace_state": "vnd1=abcd,vnd2=xyz"
  },
  "downstream": {
    "tracer_config": {
      "ActorName": "next-actor-name",
      "VendorKey": "ref",
      "TrustTraceID": true,
      "TrustSampling": true,
      "Sample": true,
      "Upsample": false
    },
    "trace": {
      "trace_id": "999ed0ff376053d0ce00566d456e9c80",
      "span_id": "4052ee7a3c350ee7",
      "parent_id": "87f1df57ec17c593",
      "sampled": true,
      "trace_parent": "00-999ed0ff376053d0ce00566d456e9c80-87f1df57ec17c593-01",
      "trace_state": "someKey=here,vnd1=abcd,vnd2=xyz"
    }
  }
}
```

### Configuration

Actors are expected to be configurable via env variables w.r.t. how they behave and participate in the trace. See `TracerConfiguration` struct definition in [api/tracer_config.go](./api/tracer_config.go). The default actor implementation supports environment variables defined at the top of that file.

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

* Implement checks in the behaviors depending on the tracer configuration.
* Use a real tracer that supports Trace Context semantics to implement the `api.Trace` and create another actor
* Perhaps merge "actor" and "node" into a single term.
