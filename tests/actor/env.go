package actor

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
	errNoConfig   = errors.New("no env configuration")
	errIncomplete = fmt.Errorf("not all env variables are defined: %v", []string{
		envTrustTraceID, envTrustSampling, envSample, envUpsample,
	})
)

func configFromEnv() (api.ActorConfiguration, error) {
	var (
		trustTraceID  = os.Getenv(envTrustTraceID)
		trustSampling = os.Getenv(envTrustSampling)
		sample        = os.Getenv(envSample)
		upsample      = os.Getenv(envUpsample)
	)
	if trustTraceID == "" && trustSampling == "" && sample == "" && upsample == "" {
		return api.ActorConfiguration{}, errNoConfig
	}
	if trustTraceID == "" || trustSampling == "" || sample == "" || upsample == "" {
		return api.ActorConfiguration{}, errIncomplete
	}
	return api.ActorConfiguration{
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
