SHELL = /bin/sh
.DEFAULT_GOAL := build
COMMIT_HASH := $(shell git rev-parse --short HEAD)

build: 
	@mkdir -p bin && \
	echo Building x-bootstrap-node for all architectures.
	GOOS=linux GOARCH=arm go build -o bin/x-bootstrap-node-linux-arm; \
	GOOS=linux GOARCH=arm64 go build -o bin/x-bootstrap-node-linux-arm64; \
	GOOS=linux GOARCH=386 go build -o bin/x-bootstrap-node-linux-386; \
	GOOS=linux GOARCH=amd64 go build -o bin/x-bootstrap-node-linux-amd64; \
