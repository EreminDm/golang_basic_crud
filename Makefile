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
	# echo "USE mysql;\nUPDATE user SET password=PASSWORD('root') WHERE user='root';\nFLUSH PRIVILEGES;\n" | mysql -u root
	sudo mysql -e "CREATE DATABASE IF NOT EXISTS person;" -u root
	sudo mysql -u root -e CREATE USER IF NOT EXISTS 'test'@'%' IDENTIFIED BY 'test';
	sudo mysql -u root -e GRANT ALL ON *.* TO 'test'@'%';
	sudo mysql -u root -e FLUSH PRIVILEGES;
	# Tweak PATH for Travis
	export PATH=$PATH:$HOME/gopath/bin
	
	set -ex

	sql-migrate status -config=db/mariadb/dbconfig.yml -env mysql
	sql-migrate up -config=db/mariadb/dbconfig.yml -env mysql
