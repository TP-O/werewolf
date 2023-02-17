#!/bin/bash
set -ex

protoc --proto_path=grpc/proto grpc/proto/*.proto --go_out=. --go-grpc_out=.
