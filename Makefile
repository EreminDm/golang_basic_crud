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
	sudo mysql -e "CREATE DATABASE IF NOT EXISTS person;" -u root
	# Tweak PATH for Travis
	export PATH=$PATH:$HOME/gopath/bin
	
	set -ex

	sql-migrate status -config=db/mariadb/dbconfig.yml -env mysql
	sql-migrate up -config=db/mariadb/dbconfig.yml -env mysql

.PHONY: installkub
installkub:
	#!/bin/sh
	curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
	chmod +x ./kubectl
	sudo mv ./kubectl /usr/local/bin/kubectl


.PHONY: preparekub
preparekub:
	if [ ! -d "$HOME/google-cloud-sdk/bin" ]; then
	rm -rf "$HOME/google-cloud-sdk"
	curl https://sdk.cloud.google.com | bash > /dev/null
	fi
	# Promote gcloud to PATH top priority (prevent using old version fromt travis)
	source $HOME/google-cloud-sdk/path.bash.inc

	# Make sure kubectl is updated to latest version
	gcloud components update kubectl

.PHONY: gauth
gauth:
		@gcloud auth activate-service-account --key-file client-secret.json
		
.PHONY: gconfig
gconfig:
		@gcloud config set project golang-basic-crud
		@gcloud container clusters \
		get-credentials gbs \
		--zone us-central1-a \
		--project golang-basic-crud
		#???
		@gcloud auth configure-docker 		
.PHONY: buildkub
buildkub:
		@docker build -t gcr.io/golang-basic-crud/gbc-f:v0.3 .
.PHONY: runkub
runkub:
		@docker run -p 8000:8000 -p 8888:8888 gcr.io/golang-basic-crud/gbc-f:v0.3
.PHONY: pushkub
pushkub:
		@docker push gcr.io/golang-basic-crud/gbc-f:v0.3		
PHONY: deploykub
deploykub: gconfig
		@kubectl apply -f k8s.yaml
		# https://github.com/kubernetes/kubernetes/issues/27081#issuecomment-238078103
		@kubectl patch deployment gbc-f -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"`date +'%s'`\"}}}}}"
		
# .PHONY: kub
# kub: 
# 	wget https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz 
# 	tar xzf google-cloud-sdk.tar.gz 
# 	./google-cloud-sdk/install.sh -q
# 	cd google-cloud-sdk
# 	gcloud -q components install kubectl
# 	#                                  projectid         cluster 
# 	gcloud builds submit --tag gcr.io/golang-basic-crud/golang_basic_crud .
# 	kubectl apply -f deployment.yaml
# 	kubectl apply -f service.yaml