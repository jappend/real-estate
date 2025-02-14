.DEFAULT_GOAL := run

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

PG_CONN_STRING=postgresql://$(PGUSER):$(PGPASSWORD)@$(PGHOST):$(PGPORT)/$(PGDATABASE)?sslmode=disable

fmt:
	gofmt -w .
.PHONY:fmt

run: fmt
		air
.PHONY:run

test:
	go test ./...
.PHONY:test

build: 
	go build -o real-estate
.PHONY:build

m-create:
ifndef name
	$(error name is required, e.g., `make m-create name=users_alter_table_add_column_email`)
endif
	goose create -s -dir internal/database/migrations $(name) sql
.PHONY: m-create
		
m-up:
	goose -dir internal/database/migrations postgres "$(PG_CONN_STRING)" up
.PHONY:m-up

m-down:
	goose -dir internal/database/migrations postgres "$(PG_CONN_STRING)" down 
.PHONY:m-up

m-status:
	goose -dir internal/database/migrations postgres "$(PG_CONN_STRING)" status 
.PHONY:m-up
