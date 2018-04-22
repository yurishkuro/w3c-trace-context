package reftracer

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/w3c/distributed-tracing/tests/api"
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
	DefaultTracerConfiguration = api.TracerConfiguration{
		ActorName:     "undefined",
		TrustTraceID:  true,
		TrustSampling: true,
		Sample:        true,
		Upsample:      false,
	}
)

// TracerConfigFromEnv reads TracerConfiguration from environment variables.
func TracerConfigFromEnv() (api.TracerConfiguration, error) {
	var (
		actorName     = os.Getenv(envActorName)
		trustTraceID  = os.Getenv(envTrustTraceID)
		trustSampling = os.Getenv(envTrustSampling)
		sample        = os.Getenv(envSample)
		upsample      = os.Getenv(envUpsample)
	)
	if actorName == "" && trustTraceID == "" && trustSampling == "" && sample == "" && upsample == "" {
		return api.TracerConfiguration{}, ErrNoConfig
	}
	if actorName == "" || trustTraceID == "" || trustSampling == "" || sample == "" || upsample == "" {
		return api.TracerConfiguration{}, ErrIncomplete
	}
	return api.TracerConfiguration{
		ActorName:     actorName,
		TrustTraceID:  toBool(trustTraceID),
		TrustSampling: toBool(trustSampling),
		Sample:        toBool(sample),
		Upsample:      toBool(upsample),
	}, nil
}

func toBool(v string) bool {
	b, err := strconv.ParseBool(v)
	if err != nil {
		panic("Cannot parse value as boolean: " + err.Error())
	}
	return b
}
