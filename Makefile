BIN_LINUX := "./bin/monitoring"
BIN_WINDOWS := "./bin/monitoring.exe"


build-linux:
	GOOS="linux" go build -tags linux -o $(BIN_LINUX) ./cmd

build-windows:
	GOOS="windows" go build -tags windows -o $(BIN_WINDOWS) ./cmd

run-linux: build-linux
	$(BIN_LINUX)

run-windows: build-windows
	$(BIN_WINDOWS)

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.52.2

lint: install-lint-deps
	golangci-lint run ./...

test: 
	go test -race ./...

generate:
	go generate ./...

.PHONY: generate lint build-linux build-windows run install-lint-deps