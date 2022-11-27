package action

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type oracle struct {
	Role    map[types.PlayerID]types.RoleID    `json:"role"`
	Faction map[types.PlayerID]types.FactionID `json:"faction"`
}

type predict struct {
	*action[*oracle]
	predictedRoleID    types.RoleID
	predictedFactionID types.FactionID
}

func NewRolePredict(game contract.Game, roleID types.RoleID) contract.Action {
	return &predict{
		action: &action[*oracle]{
			id:   config.PredictActionID,
			game: game,
			state: &oracle{
				Role: make(map[types.PlayerID]types.RoleID),
			},
		},
		predictedRoleID: roleID,
	}
}

func NewFactionPredict(game contract.Game, factionID types.FactionID) contract.Action {
	return &predict{
		action: &action[*oracle]{
			id:   config.PredictActionID,
			game: game,
			state: &oracle{
				Faction: make(map[types.PlayerID]types.FactionID),
			},
		},
		predictedFactionID: factionID,
	}
}

func (p *predict) Perform(req *types.ActionRequest) *types.ActionResponse {
	return p.action.perform(p.validate, p.execute, p.skip, req)
}

func (p *predict) validate(req *types.ActionRequest) (msg string) {
	targetID := req.TargetIDs[0]
	isKnownTarget := (p.state.Role != nil && slices.Contains(maps.Keys(p.state.Role), targetID)) ||
		(p.state.Faction != nil && slices.Contains(maps.Keys(p.state.Faction), targetID))

	if req.ActorID == targetID {
		msg = "WTF! You don't know who you are?"
	} else if isKnownTarget {
		msg = "You already knew this player :D"
	} else if player := p.game.Player(targetID); player == nil {
		msg = "Unable to see this player!"
	}

	return
}

func (p *predict) execute(req *types.ActionRequest) *types.ActionResponse {
	isCorrect := false
	target := p.game.Player(req.TargetIDs[0])

	// Check if player's faction or role is correct
	if p.state.Faction != nil {
		if target.FactionID() == p.predictedFactionID {
			p.state.Faction[target.ID()] = p.predictedFactionID
			isCorrect = true
		} else {
			p.state.Faction[target.ID()] = types.FactionID(0)
		}
	} else if p.state.Role != nil {
		if target.MainRoleID() == p.predictedRoleID {
			p.state.Role[target.ID()] = p.predictedRoleID
			isCorrect = true
		} else {
			p.state.Role[target.ID()] = types.RoleID(0)
		}
	} else {
		return &types.ActionResponse{
			Ok:      false,
			Message: "System error!",
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: isCorrect,
	}
}
