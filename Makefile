BINARY_NAME=dbk
.DEFAULT_GOAL := build

.PHONY: fmt vet test build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build:
	go build -v ./cmd/...

test:
	go test ./...

test-no-cache:
	go test -v -count=1 ./...

clean:
	go clean
	rm -f ${BINARY_NAME}

