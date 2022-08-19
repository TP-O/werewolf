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

// Get roles' name.
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

// Check condition is satisfied then if pass, activate skill
// corresponding to this role based on game context.
func (r *role) ActivateSkill(req *types.ActionRequest) *types.ActionResponse {
	if r.skill == nil ||
		r.skill.numberOfUses == types.OutOfTimes {

		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag: types.SystemErrorTag,
				Msg: map[string]string{
					types.AlertErrorField: "Unable to use skill!",
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
