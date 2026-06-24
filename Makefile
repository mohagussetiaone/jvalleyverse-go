.PHONY: build run test test-v seed clean

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./... -count=1

test-v:
	go test ./... -count=1 -v

seed:
	go run ./cmd/seed

clean:
	rm -rf bin/
