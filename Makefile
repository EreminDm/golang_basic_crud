# Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOTEST=$(GOCMD) test
	GOLANGCI=golangci-lint
all: litern test

litern:
$(GOLANGCI) run --enable-all

test:
 $(GOTEST) -v ./...
