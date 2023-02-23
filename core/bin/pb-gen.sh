#!/bin/bash
set -ex

rm -rf grpc/pb/* && \
protoc --proto_path=grpc/proto grpc/proto/*.proto --go_out=. --go-grpc_out=.
