# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLANGCI=golangci-lint

.PHONY: all
all: install lint test

.PHONY: install
install:
	$(GOCMD) get "github.com/pkg/errors"
	$(GOCMD) get "github.com/gorilla/mux" 
	$(GOCMD) get "go.mongodb.org/mongo-driver"

.PHONY: lint
lint:
	$(GOLANGCI) run

.PHONY: test
test:
	$(GOTEST) -v ./...
