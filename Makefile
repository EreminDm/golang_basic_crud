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
	$(GOCMD) get "go.mongodb.org/mongo-driver/mongo"
	$(GOCMD) get "github.com/stretchr/testify/assert"
.PHONY: lint
lint:
	$(GOLANGCI) run

.PHONY: test
test:
	$(GOTEST) -v ./...
