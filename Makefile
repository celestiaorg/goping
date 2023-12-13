BINARY_NAME := goping

build:
	go build -o bin/$(BINARY_NAME) -v ./cmd

docker:
	docker build -t goping-server .

lint:
	golangci-lint run ./...

.PHONY: build docker lint