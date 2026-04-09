.PHONY: build run-server send

build:
	go build ./...

run-server:
	go run ./cmd/server

send:
	go run ./cmd/send

kill:
	npx kill-port 8080

