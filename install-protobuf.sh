#!/bin/sh
set -ex
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protobuf-all-3.6.1.tar.gz
tar -xzvf protobuf-all-3.6.1.tar.gz
pushd protobuf-3.6.1 && ./configure --prefix=/usr && make && sudo make install && popd
