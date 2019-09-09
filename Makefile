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
	$(GOCMD) get "go.mongodb.org/mongo-driver/mongo"
	$(GOCMD) get "github.com/stretchr/testify/assert"
.PHONY: lint
lint:
	$(GOCMD) get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOLANGCI) run

.PHONY: test
test:
	$(GOTEST) -v ./...
