.PHONY: default
default: all

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build -v -ldflags "-s -w" -o dist/goello-client ./cmd/client/main.go
	go build -v -ldflags "-s -w" -o dist/goello-server ./cmd/server/main.go

.PHONY: all
all: fmt vet build
