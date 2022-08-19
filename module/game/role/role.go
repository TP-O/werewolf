package role

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

type role struct {
	id      types.RoleId
	name    string
	game    contract.Game
	player  contract.Player
	phaseId types.PhaseId
	skill   *skill
}

type skill struct {
	action       contract.Action
	numberOfUses types.NumberOfTimes
	beginRoundId types.RoundId
}

func (r *role) GetId() types.RoleId {
	return r.id
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
		r.skill.numberOfUses == types.OutOfTimes ||
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

	if res.Error != nil && r.skill.numberOfUses != types.UnlimitedTimes {
		r.skill.numberOfUses--
	}

	return res
}
