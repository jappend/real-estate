.DEFAULT_GOAL := run

fmt:
	gofmt -w .
.PHONY:fmt

run: fmt
		air
.PHONY:run
