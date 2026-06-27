.PHONY: run build tidy test

run:
	go run ./cmd/server

build:
	go build -o bin/tou-api ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...

db-shell:
	sqlite3 tou.db
