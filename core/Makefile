DB_URL="postgresql://ww_username:ww_password@postgres:5432/ww_db?sslmode=disable"

dev:
	docker-compose up

mockgen:
	mockgen -source internal/app/game/logic/contract/world.go -destination test/mock/app/game/logic/world.go -package mock_game_logic
	mockgen -source internal/app/game/logic/contract/player.go -destination test/mock/app/game/logic/player.go -package mock_game_logic
	mockgen -source internal/app/game/logic/contract/action.go -destination test/mock/app/game/logic/action.go -package mock_game_logic
	mockgen -source internal/app/game/logic/contract/role.go -destination test/mock/app/game/logic/role.go -package mock_game_logic
	mockgen -source internal/app/game/logic/contract/poll.go -destination test/mock/app/game/logic/poll.go -package mock_game_logic
	mockgen -source internal/app/game/logic/contract/scheduler.go -destination test/mock/app/game/logic/scheduler.go -package mock_game_logic
	mockgen -source internal/app/game/logic/contract/moderator.go -destination test/mock/app/game/logic/moderator.go -package mock_game_logic
	mockgen -source internal/app/game/logic/contract/provider.go -destination test/mock/app/game/logic/provider.go -package mock_game_logic
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
