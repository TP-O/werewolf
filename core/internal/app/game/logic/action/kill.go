package action

import (
	"fmt"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

// kill kills one player in the world.
type kill struct {
	action

	// Kills stores kill records. The key is player ID and
	// the value is the number of times that player was killed.
	Kills map[types.PlayerId]uint `json:"kills"`
}

func NewKill(world contract.World) contract.Action {
	return &kill{
		action: action{
			id:    KillActionId,
			world: world,
		},
		Kills: make(map[types.PlayerId]uint),
	}
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (k *kill) Execute(req types.ActionRequest) types.ActionResponse {
	return k.action.execute(k, k.Id(), &req)
}

// validate checks if the action request is valid.
func (k kill) validate(req *types.ActionRequest) error {
	if req.ActorId == req.TargetId {
		return fmt.Errorf("Appreciate your own life (｡´ ‿｀♡)")
	} else if player := k.world.Player(req.TargetId); player == nil {
		return fmt.Errorf("Player does not exist (⊙＿⊙')")
	} else if player.IsDead() {
		return fmt.Errorf("Player is dead [¬º-°]¬")
	}

	return nil
}

// perform completes the action request.
func (k *kill) perform(req *types.ActionRequest) types.ActionResponse {
	player := k.world.Player(req.TargetId)
	if player == nil {
		return types.ActionResponse{
			Ok:      false,
			Message: "The targeted player doesn't exist!",
		}
	}

	if !player.Die() {
		return types.ActionResponse{
			Ok:      false,
			Message: "Unable to kill the targeted player!",
		}
	}

	k.Kills[player.Id()]++
	return types.ActionResponse{
		Ok:   true,
		Data: player.Id(),
	}
}
