package api

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	envActorName     = "ACTOR_NAME"
	envTrustTraceID  = "TRUST_TRACE_ID"
	envTrustSampling = "TRUST_SAMPLING"
	envSample        = "SAMPLE"
	envUpsample      = "UPSAMPLE"
)

var (
	// ErrNoConfig is returned when all env variables were absent.
	ErrNoConfig = errors.New("no env configuration")

	// ErrIncomplete is returned when some env variables were absent.
	ErrIncomplete = fmt.Errorf("not all env variables are defined: %v", []string{
		envActorName, envTrustTraceID, envTrustSampling, envSample, envUpsample,
	})

	// DefaultTracerConfiguration is the most common configuration.
	DefaultTracerConfiguration = TracerConfiguration{
		ActorName:     "undefined",
		TrustTraceID:  true,
		TrustSampling: true,
		Sample:        true,
		Upsample:      false,
	}
)

// TracerConfiguration describes how the actor's tracer is going to behave under different conditions.
type TracerConfiguration struct {
	ActorName string

	// String key that this actor is using to represent its state in the `tracestate` header.
	VendorKey string

	// TrustTraceID controls whether the tracer respects inbound trace ID or creates a new trace
	// and records inbound trace ID as correlation.
	TrustTraceID bool

	// TrustSampling control whether the tracer respects inbound sampling flag or makes its own decision (based on Sample below).
	TrustSampling bool

	// Sample controls which sampling decision the tracer makes when it needs to make it (e.g when there is no inbound trace context).
	Sample bool

	// Upsample controls whether the tracer will switch on sampling even if the inbound trace context has sampling=off.
	Upsample bool
}

// TracerConfigFromEnv reads TracerConfiguration from environment variables.
func TracerConfigFromEnv() (TracerConfiguration, error) {
	var (
		actorName     = os.Getenv(envActorName)
		trustTraceID  = os.Getenv(envTrustTraceID)
		trustSampling = os.Getenv(envTrustSampling)
		sample        = os.Getenv(envSample)
		upsample      = os.Getenv(envUpsample)
	)
	if actorName == "" && trustTraceID == "" && trustSampling == "" && sample == "" && upsample == "" {
		return TracerConfiguration{}, ErrNoConfig
	}
	if actorName == "" || trustTraceID == "" || trustSampling == "" || sample == "" || upsample == "" {
		return TracerConfiguration{}, ErrIncomplete
	}

	toBool := func(v string) bool {
		b, err := strconv.ParseBool(v)
		if err != nil {
			panic("Cannot parse value as boolean: " + err.Error())
		}
		return b
	}

	return TracerConfiguration{
		ActorName:     actorName,
		TrustTraceID:  toBool(trustTraceID),
		TrustSampling: toBool(trustSampling),
		Sample:        toBool(sample),
		Upsample:      toBool(upsample),
	}, nil
}
