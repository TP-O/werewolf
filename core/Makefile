DB_URL="postgresql://ww_username:ww_password@postgres:5432/ww_db?sslmode=disable"

dev:
	docker-compose up

mockgen:
	mockgen -source game/contract/game.go -destination mock/game/game.go -package mock_game
	mockgen -source game/contract/player.go -destination mock/game/player.go -package mock_game
	mockgen -source game/contract/poll.go -destination mock/game/poll.go -package mock_game
	mockgen -source game/contract/scheduler.go -destination mock/game/scheduler.go -package mock_game
	mockgen -source game/contract/action.go -destination mock/game/action.go -package mock_game
	mockgen -source game/contract/role.go -destination mock/game/role.go -package mock_game
	mockgen -source game/contract/moderator.go -destination mock/game/moderator.go -package mock_game
	mockgen -source game/contract/manager.go -destination mock/game/manager.go -package mock_game
	mockgen -source storage/postgres/store.go -destination mock/storage/postgres.go -package mock_storage \
        -aux_files=uwwolf/storage/postgres=storage/postgres/querier.go
	mockgen -source app/service/room.go -destination mock/app/service/room.go -package mock_service
	mockgen -source app/service/game.go -destination mock/app/service/game.go -package mock_service

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
