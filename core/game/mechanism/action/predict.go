package action

import (
	"fmt"
	"uwwolf/game/declare"
	"uwwolf/game/mechanism/contract"
	"uwwolf/game/types"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// predict predicts whether a player in world belongs to a
// specific role or faction.
type predict struct {
	action

	// roleID is expected role ID.
	RoleID types.RoleID `json:"expected_role_id"`

	// factionID is expected faction ID.
	FactionID types.FactionID `json:"expected_faction_id"`

	// Role stores prediction results. The key is player ID and
	// the value is the expected role or not.
	Role map[types.PlayerID]bool `json:"role_prediction"`

	// Faction stores prediction results. The key is player ID and
	// the value is the expected faction or not.
	Faction map[types.PlayerID]bool `json:"faction_prediction"`
}

func NewRolePredict(world contract.World, roleID types.RoleID) contract.Action {
	return &predict{
		action: action{
			id:    declare.PredictActionID,
			world: world,
		},
		RoleID: roleID,
		Role:   make(map[types.PlayerID]bool),
	}
}

func NewFactionPredict(world contract.World, factionID types.FactionID) contract.Action {
	return &predict{
		action: action{
			id:    declare.PredictActionID,
			world: world,
		},
		FactionID: factionID,
		Faction:   make(map[types.PlayerID]bool),
	}
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (p *predict) Execute(req *types.ActionRequest) *types.ActionResponse {
	return p.action.execute(p, req)
}

// validate checks if the action request is valid.
func (p predict) validate(req *types.ActionRequest) error {
	isKnown := slices.Contains(maps.Keys(p.Role), req.TargetID) ||
		slices.Contains(maps.Keys(p.Faction), req.TargetID)

	if req.ActorID == req.TargetID {
		return fmt.Errorf("WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻")
	} else if isKnown {
		return fmt.Errorf("You already knew this player ¯\\(º_o)/¯")
	} else if player := p.world.Player(req.TargetID); player == nil {
		return fmt.Errorf("Non-existent player ¯\\_(ツ)_/¯")
	}

	return nil
}

// perform completes the action request.
func (p *predict) perform(req *types.ActionRequest) *types.ActionResponse {
	target := p.world.Player(req.TargetID)

	// Check if the player's faction or role is as expected
	if p.Faction != nil && !p.FactionID.IsUnknown() {
		p.Faction[target.ID()] = target.FactionID() == p.FactionID
		return &types.ActionResponse{
			Ok:   true,
			Data: p.Faction[target.ID()],
		}
	} else if p.Role != nil && !p.RoleID.IsUnknown() {
		p.Role[target.ID()] = slices.Contains(target.RoleIDs(), p.RoleID)
		return &types.ActionResponse{
			Ok:   true,
			Data: p.Role[target.ID()],
		}
	} else {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Unable to predict!",
		}
	}
}
