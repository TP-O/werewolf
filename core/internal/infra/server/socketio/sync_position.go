package socketio

import (
	"fmt"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"

	socketio "github.com/googollee/go-socket.io"
	"github.com/paulmach/orb"
)

type SyncPositionListenMessage struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type SyncPositionEmitMessage struct {
	X        float64        `json:"x"`
	Y        float64        `json:"y"`
	PlayerID types.PlayerID `json:"player_id"`
}

const syncPositionEvent = "sync_position"

func (s *Server) SyncPosition(client socketio.Conn, msg string) {
	var data SyncPositionListenMessage
	if err := util.JsonUnmarshal(msg, &data); err != nil {
		client.Emit(errorEvent, err.Error())
	}

	ctx := client.Context().(*clientContext)
	mod := s.gameManger.ModeratorOfPlayer(ctx.playerId)
	if mod != nil {
		_, err := mod.Player(ctx.playerId).Move(orb.Point{data.X, data.Y})
		if err != nil {
			client.Emit(errorEvent, err.Error())
		}

		s.BroadcastToRoom(
			defaultNamespace,
			fmt.Sprintf("%v", mod.GameID()),
			syncPositionEvent, message[SyncPositionEmitMessage]{
				Event: syncPositionEvent,
				Data: SyncPositionEmitMessage{
					X:        data.X,
					Y:        data.Y,
					PlayerID: ctx.playerId,
				},
			})
	}
}
