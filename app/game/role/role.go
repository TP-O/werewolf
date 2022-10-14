package role

import (
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type role struct {
	id           types.RoleId
	game         contract.Game
	player       contract.Player
	phaseId      types.PhaseId
	factionId    types.FactionId
	priority     int
	actions      map[types.ActionId]contract.Action
	beginRoundId types.RoundId
}

func (r *role) Id() types.RoleId {
	return r.id
}

func (r *role) PhaseId() types.PhaseId {
	return r.phaseId
}

func (r *role) FactionId() types.FactionId {
	return r.factionId
}

func (r *role) Priority() int {
	return r.priority
}

func (r *role) BeginRound() types.RoundId {
	return r.beginRoundId
}

func (r *role) Expiration() types.Expiration {
	expiration := types.Expiration(0)

	for _, action := range r.actions {
		expiration += action.Expiration()
	}

	return expiration
}

func (r *role) AfterBeingVoted() bool {
	return true
}

func (r *role) AfterDeath() {
	//
}

func (r *role) ActivateSkill(req *types.ActionRequest) *types.ActionResponse {
	for _, action := range r.actions {
		if req.ActionId == action.Id() {
			return action.Perform(req)
		}
	}

	return &types.ActionResponse{
		Ok:           false,
		PerformError: "Unable to perform this action!",
	}
}
