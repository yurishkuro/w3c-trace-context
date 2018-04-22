package params

import (
	"strconv"

	crossdock "github.com/crossdock/crossdock-go"
)

const (
	// Actor parameter defines the name of the Node that implements the actor under test.
	Actor = "actor"

	// RefActor is the name of the reference implementation of the actor.
	RefActor = "ref"

	// Server parameter is used as an override for the location of the actor.
	Server = "server"

	// BehaviorMalformedTraceContext tests how actor reacts to malformed trace context.
	BehaviorMalformedTraceContext = "malformed_trace_context"

	BehaviorMissingTraceContext = "missing_trace_context"

	BehaviorTraceContextSameVendor = "trace_context_same_vendor"

	BehaviorTraceContextDiffVendor = "trace_context_diff_vendor"
)

// GetBool returns the value of a boolean parameter, or fails if parameter not present.
func GetBool(t crossdock.T, name string) bool {
	fatals := crossdock.Fatals(t)

	val := t.Param(name)
	fatals.NotEmpty(val, "param %s must not be empty", name)

	b, err := strconv.ParseBool(val)
	fatals.NoError(err, "param %s must be true or false", name)

	return b
}
