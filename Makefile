# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLANGCI=golangci-lint

.PHONY: all
all: litern test

.PHONY: lint
lint:
$(GOLANGCI) run --enable-all

.PHONY: test
test:
 $(GOTEST) -v ./...
