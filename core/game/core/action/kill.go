package action

import (
	"fmt"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type kill struct {
	*action[types.KillState]
}

func NewKill(game contract.Game) contract.Action {
	return &kill{
		&action[types.KillState]{
			id:    enum.KillActionID,
			game:  game,
			state: make(types.KillState),
		},
	}
}

func (k *kill) Execute(req *types.ActionRequest) *types.ActionResponse {
	return k.action.combine(k.Skip, k.Validate, k.Perform, req)
}

func (k *kill) Validate(req *types.ActionRequest) error {
	if err := k.action.Validate(req); err != nil {
		return err
	}

	targetID := req.TargetIDs[0]

	if req.ActorID == targetID {
		return fmt.Errorf("Appreciate your own life (｡´ ‿｀♡)")
	} else if player := k.game.Player(targetID); player == nil {
		return fmt.Errorf("Player does not exist (⊙＿⊙')")
	} else if player.IsDead() {
		return fmt.Errorf("Player is dead [¬º-°]¬")
	}

	return nil
}

func (k *kill) Perform(req *types.ActionRequest) *types.ActionResponse {
	killedPlayer := k.game.KillPlayer(req.TargetIDs[0], false)
	k.state[killedPlayer.ID()]++

	return &types.ActionResponse{
		Ok:   true,
		Data: killedPlayer.ID(),
	}
}
