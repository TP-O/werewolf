#!/bin/sh
set -ex

mockgen -source module/game/contract/game.go  -destination mock/game/game.go -package game
mockgen -source module/game/contract/player.go  -destination mock/game/player.go -package game
mockgen -source module/game/contract/role.go  -destination mock/game/role.go -package game
