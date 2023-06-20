VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
LDFLAGS = -X main.version=${VERSION}

build:
	go build -ldflags '$(LDFLAGS)' cmd/proxmox-interacter.go

install:
	go install -ldflags '$(LDFLAGS)' cmd/proxmox-interacter.go

lint:
	golangci-lint run --fix ./...
