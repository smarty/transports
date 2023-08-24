#!/usr/bin/make -f

test:
	GORACE="atexit_sleep_ms=50" go test -timeout=1s -race -covermode=atomic ./...

compile:
	go build ./...

build: test compile

.PHONY: test compile build
