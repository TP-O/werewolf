#!/bin/sh
set -ex

go run ./database/migration/migrate/main.go
