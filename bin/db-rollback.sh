#!/bin/sh
set -ex

go run ./database/migration/rollback/main.go
