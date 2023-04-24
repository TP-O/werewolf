DB_URL="postgresql://ww_username:ww_password@postgres:5432/ww_db?sslmode=disable"

dev:
	docker-compose up

mockgen:
	mockgen -source game/mechanism/contract/world.go -destination util/mock/game/world.go -package mock_game
	mockgen -source game/mechanism/contract/player.go -destination util/mock/game/player.go -package mock_game
	mockgen -source game/mechanism/contract/action.go -destination util/mock/game/action.go -package mock_game
	mockgen -source game/mechanism/contract/role.go -destination util/mock/game/role.go -package mock_game
	mockgen -source game/tool/poll.go -destination util/mock/game/poll.go -package mock_game
	mockgen -source game/tool/scheduler.go -destination util/mock/game/scheduler.go -package mock_game
	mockgen -source game/moderator.go -destination util/mock/game/moderator.go -package mock_game
	mockgen -source game/manager.go -destination util/mock/game/manager.go -package mock_game
	mockgen -source db/store.go -destination util/mock/db/postgres.go -package mock_storage \
        -aux_files=uwwolf/db=db/querier.go
	mockgen -source server/service/room.go -destination util/mock/server/service/room.go -package mock_service
	mockgen -source server/service/game.go -destination util/mock/server/service/game.go -package mock_service

protobuf:
	rm -rf grpc/pb/*
	protoc --proto_path=grpc/proto grpc/proto/*.proto --go_out=. --go-grpc_out=.

db-migrate:
	docker-compose exec app migrate -path db/migration \
        -verbose \
		-database ${DB_URL} up

db-migrate-1:
	docker-compose exec app migrate -path db/migration \
        -verbose \
		-database ${DB_URL} up 1

db-rollback:
	docker-compose exec app migrate -path db/migration \
        -verbose \
		-database ${DB_URL} down

db-rollback-1:
	docker-compose exec app migrate -path db/migration \
        -verbose \
		-database ${DB_URL} down 1

db-drop:
	docker-compose exec app migrate -path db/migration \
        -verbose \
		-database ${DB_URL} drop

db-seed:
	docker-compose exec app echo "seed"

db-sqlc:
	docker-compose exec app sqlc generate
