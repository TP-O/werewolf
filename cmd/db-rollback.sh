#!/bin/sh
set -ex

go run ./db/migration/rollback/main.go
