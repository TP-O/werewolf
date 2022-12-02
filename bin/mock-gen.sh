#!/bin/bash
set -ex

mockgen -source app/game/contract/game.go  -destination mock/game/game.go -package gamemock
mockgen -source app/game/contract/player.go  -destination mock/game/player.go -package gamemock
mockgen -source app/game/contract/poll.go  -destination mock/game/poll.go -package gamemock
mockgen -source app/game/contract/scheduler.go  -destination mock/game/scheduler.go -package gamemock
mockgen -source app/game/contract/action.go  -destination mock/game/action.go -package gamemock
