package service

import (
	"context"
	"encoding/json"
	"fmt"
	"uwwolf/game/types"
	"uwwolf/server/data"
	"uwwolf/server/enum"

	"github.com/redis/go-redis/v9"
)

// RommService handles room-related business logic.
type RoomService interface {
	// PlayerWaitingRoom returns the room containing the given player ID.
	PlayerWaitingRoom(playerID types.PlayerID) (data.WaitingRoom, bool)
}

type roomService struct {
	// rdb is redis connection.
	rdb *redis.ClusterClient
}

func NewRoomService(rdb *redis.ClusterClient) RoomService {
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
`, enum.JoinedRoomIdNs, enum.RoomNs)

// PlayerWaitingRoom returns the room containing the given player ID.
func (rs roomService) PlayerWaitingRoom(playerID types.PlayerID) (data.WaitingRoom, bool) {
	var room data.WaitingRoom

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
