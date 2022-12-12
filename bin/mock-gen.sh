#!/bin/bash
set -ex

mockgen -source game/contract/game.go  -destination mock/game/game.go -package gamemock
mockgen -source game/contract/player.go  -destination mock/game/player.go -package gamemock
mockgen -source game/contract/poll.go  -destination mock/game/poll.go -package gamemock
mockgen -source game/contract/scheduler.go  -destination mock/game/scheduler.go -package gamemock
mockgen -source game/contract/action.go  -destination mock/game/action.go -package gamemock
mockgen -source game/contract/role.go  -destination mock/game/role.go -package gamemock
