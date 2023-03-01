dev:
	docker-compose up

mockgen:
	mockgen -source game/contract/game.go  -destination mock/game/game.go -package gamemock
	mockgen -source game/contract/player.go  -destination mock/game/player.go -package gamemock
	mockgen -source game/contract/poll.go  -destination mock/game/poll.go -package gamemock
	mockgen -source game/contract/scheduler.go  -destination mock/game/scheduler.go -package gamemock
	mockgen -source game/contract/action.go  -destination mock/game/action.go -package gamemock
	mockgen -source game/contract/role.go  -destination mock/game/role.go -package gamemock

protobuf:
	rm -rf grpc/pb/*
	protoc --proto_path=grpc/proto grpc/proto/*.proto --go_out=. --go-grpc_out=.

db/migrate:
	docker-compose exec app go run ./db/migration/main.go up

db/drop:
	docker-compose exec app go run ./db/migration/main.go drop

db/seed:
	docker-compose exec app go run ./db/seeder/main.go
