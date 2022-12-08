#!/bin/sh
set -ex

go run ./db/migration/main.go $@
