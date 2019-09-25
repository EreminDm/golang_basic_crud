# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLANGCI=golangci-lint

.PHONY: all
all: proto test lint mysql

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: lint
lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.18.0
	$(GOLANGCI) run

.PHONY: proto
proto:
	sudo wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protobuf-all-3.6.1.tar.gz
	sudo cp -R protobuf-all-3.6.1.tar.gz ~/go/src/github.com
	sudo rm -rf protobuf-all-3.6.1.tar.gz
	cd ~/go/src/github.com
	sudo tar -xzvf protobuf-all-3.6.1.tar.gz
	sudo pushd protobuf-3.6.1 && ./configure --prefix=/usr && make && sudo make install && popd
	cd ~/go/src/github.com/EreminDm/golang_basic_crud/net/grpc
	sudo protoc -I ./net/grpc/proto grpc.proto --go_out=plugins=grpc:./net/grpc/proto

.PHONY: mysql
mysql:
	#!/bin/bash

	# Tweak PATH for Travis
	export PATH=$PATH:$HOME/gopath/bin
	export MYSQL_USER=root
	export DATABASE_NAME=test_env
	export MYSQL_PASSWORD=test
	export MYSQL_HOST=localhost
	export MYSQL_PORT=3306

	OPTIONS="-config=db/mariadb/dbconfig.yml -env mysql"

	set -ex

	sql-migrate status $OPTIONS
	sql-migrate up $OPTIONS
