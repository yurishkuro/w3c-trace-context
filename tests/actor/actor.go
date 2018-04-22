package actor

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/w3c/distributed-tracing/tests/api"
)

// Actor implements an actor in the test suite, given a vendor-specific tracer
type Actor struct {
	server *http.Server
	tracer api.Tracer
	config api.ActorConfiguration
}

// New creates a new actor.
func New(tracer api.Tracer) *Actor {
	config, err := configFromEnv()
	if err != errNoConfig {
		panic(err.Error())
	}
	if err == errNoConfig {
		config = api.ActorConfiguration{
			TrustTraceID:  true,
			TrustSampling: true,
			Sample:        true,
			Upsample:      false,
		}
	}
	return NewWithConfig(tracer, config)
}

// NewWithConfig creates a new actor with given configuration.
func NewWithConfig(tracer api.Tracer, config api.ActorConfiguration) *Actor {
	a := &Actor{
		tracer: tracer,
		config: config,
	}

	m := http.NewServeMux()
	m.HandleFunc("/trace", a.Trace)

	a.server = &http.Server{Addr: ":8081", Handler: m}
	return a
}

// Start registers actor endpoints and starts the server(s).
func (a *Actor) Start() {
	log.Print("starting actor")
	go a.serve()
	log.Print("actor started")
}

func (a *Actor) serve() {
	if err := a.server.ListenAndServe(); err != nil && !strings.Contains(err.Error(), "Server closed") {
		log.Fatalf("actor server failed: %s", err)
	}
}

// Stop shuts down the servers
func (a *Actor) Stop() {
	a.server.Shutdown(context.Background())
}
