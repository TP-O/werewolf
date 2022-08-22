package role

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

type role struct {
	id        types.RoleId
	name      string
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

func (r *role) GetId() types.RoleId {
	return r.id
}

func (r *role) GetFactionId() types.FactionId {
	return r.factionId
}

func (r *role) GetName() string {
	return r.name
}

// Do something after being voted. Return false if exonerated,
// otherwise return true.
func (r *role) AfterBeingVoted() bool {
	return true
}

// Do something before death
func (r *role) AfterDeath() {
	//
}

// Check condition is satisfied then if pass, perform action
// corresponding to this role.
func (r *role) ActivateSkill(req *types.ActionRequest) *types.ActionResponse {
	if r.skill == nil ||
		r.skill.expiration == types.OutOfTimes ||
		r.skill.beginRoundId < r.game.GetCurrentRoundId() ||
		r.phaseId != r.game.GetCurrentPhaseId() {

		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag: types.UnauthorizedErrorTag,
				Msg: map[string]string{
					types.AlertErrorField: "Unable to execute action!",
				},
			},
		}
	}

	res := r.skill.action.Execute(req)

	if res.Error != nil && r.skill.expiration != types.UnlimitedTimes {
		r.skill.expiration--
	}

	return res
}
