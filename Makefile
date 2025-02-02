.DEFAULT_GOAL := run

fmt:
	gofmt -w .
.PHONY:fmt

run: fmt
		air
.PHONY:run

test:
	go test ./...
.PHONY:test

build: fmt
	go build -o real-estate
.PHONY:build
