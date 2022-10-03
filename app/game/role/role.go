package role

import (
	"uwwolf/app/enum"
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
	score        int
	set          int
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

func (r *role) Score() int {
	return r.score
}

func (r *role) Set() int {
	return r.set
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
	if r.beginRoundId < r.game.Round().CurrentId() ||
		r.phaseId != r.game.Round().CurrentPhaseId() {

		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   enum.InvalidInputErrorTag,
				Alert: "Please wait for your turn!",
			},
		}
	}

	for _, action := range r.actions {
		if req.ActionId == action.Id() {
			return action.Perform(req)
		}
	}

	return &types.ActionResponse{
		Error: &types.ErrorDetail{
			Tag:   enum.ForbiddenErrorTag,
			Alert: "Unable to perform this action!",
		},
	}
}
