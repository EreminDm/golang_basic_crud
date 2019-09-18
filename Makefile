# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLANGCI=golangci-lint

.PHONY: all
all: fullproto test lint  

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: lint
lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.18.0
	$(GOLANGCI) run

.PHONY: fullproto
fullproto:
	go get -u github.com/golang/protobuf/protoc-gen-go
	cd ./net/grpc/; \
	protoc -I . grpc.proto --go_out=plugins=grpc:.