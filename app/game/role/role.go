package role

import (
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type role struct {
	id        types.RoleId
	game      contract.Game
	player    contract.Player
	phaseId   types.PhaseId
	factionId types.FactionId
	skill     *skill
}

type skill struct {
	action       contract.Action
	beginRoundId types.RoundId
	expiration   types.NumberOfTimes
}

func (r *role) Id() types.RoleId {
	return r.id
}

func (r *role) FactionId() types.FactionId {
	return r.factionId
}

func (r *role) AfterBeingVoted() bool {
	return true
}

func (r *role) AfterDeath() {
	//
}

func (r *role) ActivateSkill(req *types.ActionRequest) *types.ActionResponse {
	if r.skill == nil ||
		r.skill.expiration == types.OutOfTimes ||
		r.skill.beginRoundId < r.game.Round().CurrentId() ||
		r.phaseId != r.game.Round().CurrentPhaseId() {

		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   types.UnauthorizedErrorTag,
				Alert: "Unable to execute action!",
			},
		}
	}

	res := r.skill.action.Perform(req)

	if res.Error != nil && r.skill.expiration != types.UnlimitedTimes {
		r.skill.expiration--
	}

	return res
}
