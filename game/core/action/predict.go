package action

import (
	"errors"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type predict struct {
	*action[*types.PredictState]
	predictedRoleID    enum.RoleID
	predictedFactionID enum.FactionID
}

func NewRolePredict(game contract.Game, roleID enum.RoleID) contract.Action {
	return &predict{
		action: &action[*types.PredictState]{
			id:   enum.PredictActionID,
			game: game,
			state: &types.PredictState{
				Role: make(map[enum.PlayerID]enum.RoleID),
			},
		},
		predictedRoleID: roleID,
	}
}

func NewFactionPredict(game contract.Game, factionID enum.FactionID) contract.Action {
	return &predict{
		action: &action[*types.PredictState]{
			id:   enum.PredictActionID,
			game: game,
			state: &types.PredictState{
				Faction: make(map[enum.PlayerID]enum.FactionID),
			},
		},
		predictedFactionID: factionID,
	}
}

func (p *predict) Execute(req *types.ActionRequest) *types.ActionResponse {
	return p.action.combine(p.Skip, p.Validate, p.Perform, req)
}

func (p *predict) Validate(req *types.ActionRequest) error {
	if err := p.action.Validate(req); err != nil {
		return err
	}

	targetID := req.TargetIDs[0]
	isKnownTarget := (p.state.Role != nil && slices.Contains(maps.Keys(p.state.Role), targetID)) ||
		(p.state.Faction != nil && slices.Contains(maps.Keys(p.state.Faction), targetID))

	if req.ActorID == targetID {
		return errors.New("WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻")
	} else if isKnownTarget {
		return errors.New("You already knew this player ¯\\(º_o)/¯")
	} else if player := p.game.Player(targetID); player == nil {
		return errors.New("Non-existent player ¯\\_(ツ)_/¯")
	}

	return nil
}

func (p *predict) Perform(req *types.ActionRequest) *types.ActionResponse {
	isCorrect := false
	target := p.game.Player(req.TargetIDs[0])

	// Check if player's faction or role is correct
	if p.state.Faction != nil {
		if target.FactionID() == p.predictedFactionID {
			p.state.Faction[target.ID()] = p.predictedFactionID
			isCorrect = true
		} else {
			p.state.Faction[target.ID()] = enum.FactionID(0)
		}
	} else if p.state.Role != nil {
		if slices.Contains(target.RoleIDs(), p.predictedRoleID) {
			p.state.Role[target.ID()] = p.predictedRoleID
			isCorrect = true
		} else {
			p.state.Role[target.ID()] = enum.RoleID(0)
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
