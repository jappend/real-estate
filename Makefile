.DEFAULT_GOAL := run

fmt:
	gofmt -w .
.PHONY:fmt

run: fmt
		air
.PHONY:run

build: fmt
	go build -o real-estate
.PHONY:build
