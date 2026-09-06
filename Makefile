.PHONY: build test lint run postgres-up postgres-down

build:
	mkdir -p bin
	go build -o bin/migrationlab ./cmd/migrationlab

test:
	go test ./...

lint:
	@test -z "$$(gofmt -l .)" || (echo "Go files need formatting; run gofmt -w ." && gofmt -l . && exit 1)
	go vet ./...

run:
	go run ./cmd/migrationlab --help

postgres-up:
	docker compose up -d postgres

postgres-down:
	docker compose down
