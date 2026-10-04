.PHONY: run build test lint tidy clean

run:
	go run ./cmd/monitor

build:
	go build -o bin/monitor ./cmd/monitor

test:
	go test ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

clean:
	rm -rf bin/
	rm -f cloudmonitor.db
