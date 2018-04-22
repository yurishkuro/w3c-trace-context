package main

import (
	"net/url"
	"testing"

	crossdock "github.com/crossdock/crossdock-go"
	"github.com/w3c/distributed-tracing/tests/actor"
	"github.com/w3c/distributed-tracing/tests/driver"
	driverParams "github.com/w3c/distributed-tracing/tests/driver/params"
	"github.com/w3c/distributed-tracing/tests/internal/reftracer"
)

const clientURL = "http://127.0.0.1:8080"

func TestCrossdock(t *testing.T) {
	actor := actor.New(reftracer.New())
	actor.Start()
	defer actor.Stop()
	go driver.Start()

	crossdock.Wait(t, clientURL, 10)

	type params map[string]string
	type axes map[string][]string

	defaultParams := params{"server": "127.0.0.1"}

	behaviors := []struct {
		name   string
		params params
		axes   axes
	}{
		{
			name: driverParams.BehaviorMalformedTraceContext,
			axes: axes{
				"actor": []string{driverParams.RefActor},
			},
			params: params{},
		},
		{
			name: driverParams.BehaviorMissingTraceContext,
			axes: axes{
				"actor": []string{driverParams.RefActor},
			},
			params: params{},
		},
		{
			name: driverParams.BehaviorTraceContextSameVendor,
			axes: axes{
				"actor": []string{driverParams.RefActor},
			},
			params: params{},
		},
		{
			name: driverParams.BehaviorTraceContextDiffVendor,
			axes: axes{
				"actor": []string{driverParams.RefActor},
			},
			params: params{},
		},
	}

	for _, bb := range behaviors {
		t.Run(bb.name, func(t *testing.T) {
			args := url.Values{}
			for k, v := range defaultParams {
				args.Set(k, v)
			}
			for k, v := range bb.params {
				args.Set(k, v)
			}

			if len(bb.axes) == 0 {
				crossdock.Call(t, clientURL, bb.name, args)
				return
			}

			for _, entry := range crossdock.Combinations(bb.axes) {
				entryArgs := url.Values{}
				for k := range args {
					entryArgs.Set(k, args.Get(k))
				}
				for k, v := range entry {
					entryArgs.Set(k, v)
				}

				crossdock.Call(t, clientURL, bb.name, entryArgs)
			}
		})
	}
}
