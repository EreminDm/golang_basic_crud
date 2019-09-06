# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLANGCI=golangci-lint

.PHONY: all
all: lint test

.PHONY: lint
lint:
	$(GOLANGCI) run

.PHONY: test
test:
	$(GOTEST) -v ./...