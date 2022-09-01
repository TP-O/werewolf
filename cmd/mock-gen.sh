#!/bin/sh
set -ex

mockgen -source module/game/contract/game.go  -destination test/mock/game/game.go -package game
mockgen -source module/game/contract/player.go  -destination test/mock/game/player.go -package game
mockgen -source module/game/state/round.go  -destination test/mock/game/state/round.go -package state_test
