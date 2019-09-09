# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLANGCI=golangci-lint

.PHONY: all
all: install mlint test

.PHONY: install
install:
	$(GOCMD) get "github.com/pkg/errors"
	$(GOCMD) get "github.com/gorilla/mux" 
	dep ensure -add "go.mongodb.org/mongo-driver/mongo@~1.1.0"

.PHONY: lint
lint:
	$(GOLANGCI) run

.PHONY: test
test:
	$(GOTEST) -v ./...
