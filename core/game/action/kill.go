package action

import (
	"fmt"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

// kill kills one player in the game.
type kill struct {
	action

	// Kills stores kill records. The key is player ID and
	// the value is the number of times that player was killed.
	Kills map[types.PlayerID]uint `json:"kills"`
}

func NewKill(game contract.Game) contract.Action {
	return &kill{
		action: action{
			id:   KillActionID,
			game: game,
		},
		Kills: make(map[types.PlayerID]uint),
	}
}

func (k *kill) Execute(req types.ActionRequest) types.ActionResponse {
	return k.action.execute(k, req)
}

func (k kill) validate(req types.ActionRequest) error {
	if req.ActorID == req.TargetID {
		return fmt.Errorf("Appreciate your own life (｡´ ‿｀♡)")
	} else if player := k.game.Player(req.TargetID); player == nil {
		return fmt.Errorf("Player does not exist (⊙＿⊙')")
	} else if player.IsDead() {
		return fmt.Errorf("Player is dead [¬º-°]¬")
	}

	return nil
}

func (k *kill) perform(req types.ActionRequest) types.ActionResponse {
	killedPlayer := k.game.KillPlayer(req.TargetID)
	k.Kills[killedPlayer.ID()]++

	return types.ActionResponse{
		Ok: true,
		StateChanges: types.StateChanges{
			DeadPlayerID: killedPlayer.ID(),
		},
	}
}
