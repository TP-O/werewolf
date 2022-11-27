package action

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
)

type kill struct {
	*action[[]types.PlayerID]
}

func NewKill(game contract.Game) contract.Action {
	return &kill{
		&action[[]types.PlayerID]{
			id:    config.KillActionID,
			game:  game,
			state: make([]types.PlayerID, 0),
		},
	}
}

func (k *kill) Perform(req *types.ActionRequest) *types.ActionResponse {
	return k.action.perform(k.validate, k.execute, k.skip, req)
}

func (k *kill) validate(req *types.ActionRequest) (msg string) {
	targetID := req.TargetIDs[0]

	if req.ActorID == targetID {
		msg = "Appreciate your own life <3"
	} else if player := k.game.Player(targetID); player == nil || player.IsDead() {
		msg = "Unable to kill this player!"
	}

	return
}

func (k *kill) execute(req *types.ActionRequest) *types.ActionResponse {
	killedPlayer := k.game.KillPlayer(req.TargetIDs[0])
	k.state = append(k.state, killedPlayer.ID())

	return &types.ActionResponse{
		Ok:   true,
		Data: killedPlayer.ID(),
	}
}
