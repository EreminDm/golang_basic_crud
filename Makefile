# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLANGCI=golangci-lint

.PHONY: all
all: test lint 

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: lint
lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.17.0
	$(GOLANGCI) run


