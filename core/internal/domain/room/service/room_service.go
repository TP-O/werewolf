package service

import (
	"context"
	"encoding/json"
	"fmt"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/internal/domain/room/model"
	"uwwolf/internal/infra/db/redis"

	goredis "github.com/redis/go-redis/v9"
)

// RommService handles room-related business logic.
type RoomService interface {
	// PlayerWaitingRoom returns the room containing the given player ID.
	PlayerWaitingRoom(playerID types.PlayerID) (model.WaitingRoom, bool)
}

type roomService struct {
	// rdb is redis connection.
	rdb *goredis.ClusterClient
}

func NewRoomService(rdb *goredis.ClusterClient) RoomService {
	return &roomService{
		rdb,
	}
}

var GetWaitingRoomScript = fmt.Sprintf(`
    local player_id = ARGV[1]
    local room_id = redis.call("GET", %q + player_id)
    if not room_id then
        return nil
    end

    return redis.call("GET", %q + room_id)
`, redis.GameSetting, redis.WaitingRoom)

// PlayerWaitingRoom returns the room containing the given player ID.
func (rs roomService) PlayerWaitingRoom(playerID types.PlayerID) (model.WaitingRoom, bool) {
	var room model.WaitingRoom

	encodedRoom := rs.rdb.Eval(
		context.Background(),
		GetWaitingRoomScript,
		[]string{},
		playerID,
	).Val()
	if err := json.Unmarshal([]byte(fmt.Sprintf("%v", encodedRoom)), &room); err != nil {
		return room, false
	}

	return room, true
}
