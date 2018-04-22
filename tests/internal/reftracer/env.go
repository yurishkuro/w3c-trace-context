package reftracer

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/w3c/distributed-tracing/tests/api"
)

const (
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
		envTrustTraceID, envTrustSampling, envSample, envUpsample,
	})

	// DefaultTracerConfiguration is the most common configuration.
	DefaultTracerConfiguration = api.TracerConfiguration{
		TrustTraceID:  true,
		TrustSampling: true,
		Sample:        true,
		Upsample:      false,
	}
)

// TracerConfigFromEnv reads TracerConfiguration from environment variables.
func TracerConfigFromEnv() (api.TracerConfiguration, error) {
	var (
		trustTraceID  = os.Getenv(envTrustTraceID)
		trustSampling = os.Getenv(envTrustSampling)
		sample        = os.Getenv(envSample)
		upsample      = os.Getenv(envUpsample)
	)
	if trustTraceID == "" && trustSampling == "" && sample == "" && upsample == "" {
		return api.TracerConfiguration{}, ErrNoConfig
	}
	if trustTraceID == "" || trustSampling == "" || sample == "" || upsample == "" {
		return api.TracerConfiguration{}, ErrIncomplete
	}
	return api.TracerConfiguration{
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
