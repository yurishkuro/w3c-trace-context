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
	name   string
	server *http.Server
	tracer api.Tracer
}

// New creates a new actor.
func New(tracer api.Tracer) *Actor {
	a := &Actor{
		name:   tracer.Configuration().ActorName,
		tracer: tracer,
	}

	m := http.NewServeMux()
	m.HandleFunc("/trace", a.Trace)

	a.server = &http.Server{Addr: ":8081", Handler: m}
	return a
}

// Start registers actor endpoints and starts the server(s).
func (a *Actor) Start() {
	log.Printf("starting actor '%s'", a.name)
	go a.serve()
	log.Printf("actor '%s' started", a.name)
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
