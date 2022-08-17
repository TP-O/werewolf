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
	skills  []*skill
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
func (r *role) ActivateSkill(data *types.ActionData) *types.PerformResult {
	for _, skill := range r.skills {
		if skill.numberOfUses != types.OutOfTimes &&
			skill.beginRoundId >= r.game.GetCurrentRoundId() {

			res := skill.action.Execute(data)

			if res.Errors != nil && skill.numberOfUses != types.UnlimitedTimes {
				skill.numberOfUses--
			}

			return res
		}
	}

	return &types.PerformResult{
		ErrorTag: types.SystemErrorTag,
		Errors: map[string]string{
			types.SystemErrorProperty: "Skill is unavailable!",
		},
	}
}
