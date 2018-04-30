TESTBED_YAML=docker-compose.yaml

test:
	go test -v ./

crossdock-binary:
	GOOS=linux CGO_ENABLED=0 installsuffix=cgo go build -o ./crossdock-linux ./main.go

example-binary:
	GOOS=linux CGO_ENABLED=0 installsuffix=cgo go build -o ./example/example-linux ./example/main.go

crossdock: crossdock-binary
	docker-compose -f $(TESTBED_YAML) kill
	docker-compose -f $(TESTBED_YAML) rm -f refnode
	$(MAKE) crossdock-run

crossdock-fresh: crossdock-binary
	docker-compose -f $(TESTBED_YAML) kill
	docker-compose -f $(TESTBED_YAML) rm --force
	docker-compose -f $(TESTBED_YAML) pull
	docker-compose -f $(TESTBED_YAML) build
	$(MAKE) crossdock-run

crossdock-run:
	docker-compose -f $(TESTBED_YAML) run crossdock 2>&1 | tee run-crossdock.log
	grep 'Tests passed!' run-crossdock.log || (make crossdock-logs; echo TESTS FAILED; false)

crossdock-logs:
	docker-compose -f $(TESTBED_YAML) logs

crossdock-clean:
	docker-compose -f $(TESTBED_YAML) down

.PHONY: crossdock crossdock-binary crossdock-fresh crossdock-logs crossdock-clean crossdock-run
