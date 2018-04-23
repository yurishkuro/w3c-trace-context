package reftracer

import (
	"log"

	"github.com/w3c/distributed-tracing/tests/api"
	"github.com/w3c/distributed-tracing/tests/internal/random"
)

// Tracer is a reference implementation of api.Tracer
type Tracer struct {
	config api.TracerConfiguration
}

// New creates a new reference Tracer.
func New() *Tracer {
	config, err := api.TracerConfigFromEnv()
	if err == api.ErrNoConfig {
		config = api.DefaultTracerConfiguration
	} else if err != nil {
		panic(err.Error())
	}
	return NewWithConfig(config)
}

// NewWithConfig creates a new reference Tracer with given configuration.
func NewWithConfig(config api.TracerConfiguration) *Tracer {
	cfg := config // copy
	cfg.VendorKey = "ref"
	return &Tracer{
		config: cfg,
	}
}

// StartSpan implements Tracer API.
func (t *Tracer) StartSpan(tc api.TraceContext) api.Span {
	// TODO tc.TraceState should take priority
	traceID, parentSpanID, sampled, _ := tc.ParseTraceParent()
	correlationID := ""
	if traceID == "" {
		traceID = random.New64BitID() + random.New64BitID()
		sampled = t.config.Sample
	} else {
		if !t.config.TrustTraceID {
			correlationID = traceID
			log.Printf("captured correlationID=%s", correlationID)
			traceID = random.New64BitID() + random.New64BitID()
			log.Printf("restarting trace with traceID=%s", traceID)
		}
		if !t.config.TrustSampling {
			sampled = t.config.Sample
		} else {
			if !sampled && t.config.Upsample {
				sampled = t.config.Sample
			}
		}
	}
	return &Span{
		traceID:       traceID,
		spanID:        random.New64BitID(),
		parentSpanID:  parentSpanID,
		correlationID: correlationID, // TODO should depend on the participation mode
		sampled:       sampled,
		traceState:    tc.TraceState,
	}
}

// Configuration implements Tracer API.
func (t *Tracer) Configuration() api.TracerConfiguration {
	return t.config
}
