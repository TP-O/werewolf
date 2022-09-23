#!/bin/sh
set -ex

go run ./db/migration/migrate/main.go
