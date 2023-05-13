package action

import (
	"errors"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// predict predicts whether a player in world belongs to a
// specific role or faction.
type predict struct {
	action

	// roleId is expected role ID.
	RoleId types.RoleId `json:"expected_role_id"`

	// factionId is expected faction ID.
	FactionId types.FactionId `json:"expected_faction_id"`

	// Role stores prediction results. The key is player ID and
	// the value is the expected role or not.
	Role map[types.PlayerId]bool `json:"role_prediction"`

	// Faction stores prediction results. The key is player ID and
	// the value is the expected faction or not.
	Faction map[types.PlayerId]bool `json:"faction_prediction"`
}

func NewRolePredict(world contract.World, roleId types.RoleId) contract.Action {
	return &predict{
		action: action{
			id:    constants.PredictActionId,
			world: world,
		},
		RoleId: roleId,
		Role:   make(map[types.PlayerId]bool),
	}
}

func NewFactionPredict(world contract.World, factionId types.FactionId) contract.Action {
	return &predict{
		action: action{
			id:    constants.PredictActionId,
			world: world,
		},
		FactionId: factionId,
		Faction:   make(map[types.PlayerId]bool),
	}
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (p *predict) Execute(req types.ActionRequest) types.ActionResponse {
	return p.action.execute(p, p.Id(), &req)
}

// validate checks if the action request is valid.
func (p predict) validate(req *types.ActionRequest) error {
	isKnown := slices.Contains(maps.Keys(p.Role), req.TargetId) ||
		slices.Contains(maps.Keys(p.Faction), req.TargetId)

	if req.ActorId == req.TargetId {
		return errors.New("WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻")
	} else if isKnown {
		return errors.New("You already knew this player ¯\\(º_o)/¯")
	} else if player := p.world.Player(req.TargetId); player == nil {
		return errors.New("Non-existent player ¯\\_(ツ)_/¯")
	}

	return nil
}

// perform completes the action request.
func (p *predict) perform(req *types.ActionRequest) types.ActionResponse {
	target := p.world.Player(req.TargetId)

	// Check if the player's faction or role is as expected
	if p.Faction != nil && !util.IsZero(p.FactionId) {
		p.Faction[target.Id()] = target.FactionId() == p.FactionId
		return types.ActionResponse{
			Ok:   true,
			Data: p.Faction[target.Id()],
		}
	} else if p.Role != nil && !util.IsZero(p.RoleId) {
		p.Role[target.Id()] = slices.Contains(target.RoleIds(), p.RoleId)
		return types.ActionResponse{
			Ok:   true,
			Data: p.Role[target.Id()],
		}
	} else {
		return types.ActionResponse{
			Ok:      false,
			Message: "Unable to predict!",
		}
	}
}
