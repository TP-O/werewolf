package service

import (
	"context"
	"encoding/json"
	"fmt"
	"uwwolf/app/data"
	"uwwolf/game/types"

	"github.com/redis/go-redis/v9"
)

type RoomService interface {
	PlayerWaitingRoom(playerID types.PlayerID) *data.WaitingRoom
}

type roomService struct {
	rdb *redis.ClusterClient
}

func NewRoomService(rdb *redis.ClusterClient) RoomService {
	return &roomService{
		rdb,
	}
}

const (
	PlayerID2RoomIDRedisNamespace = "pID-rID:"
	WaitingRoomRedisNamespace     = "waiting-room:"
)

var queryRoom = redis.NewScript(fmt.Sprintf(`
    local player_id = ARGV[1]
    local room_id = redis.call("GET", %q + player_id)
    if not room_id then
        return nil
    end

    return redis.call("GET", %q + room_id)
`, PlayerID2RoomIDRedisNamespace, WaitingRoomRedisNamespace))

// PlayerWaitingRoom returns the room containing the given player ID.
func (rs roomService) PlayerWaitingRoom(playerID types.PlayerID) *data.WaitingRoom {
	var room *data.WaitingRoom
	roomJson := queryRoom.Run(context.Background(), rs.rdb, []string{}, playerID).String()
	if err := json.Unmarshal([]byte(roomJson), room); err != nil {
		return nil
	}

	return room
}
