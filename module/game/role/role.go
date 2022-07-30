package role

import (
	"errors"

	"uwwolf/module/game/action"
	"uwwolf/module/game/core"
	"uwwolf/types"
)

type role struct {
	id           types.RoleId
	name         string
	game         core.Game
	phaseId      types.PhaseId
	activeSkill  *activeSkill
	passiveSkill *passiveSkill
}

type Role interface {
	// Get roles' name.
	GetName() string

	// Check if role has after-death passive skill.
	HasAfterDeathSkill() bool

	// Check if role has before-death passive skill.
	HasBeforeDeathSkill() bool

	// Check condition is satisfied then perform an action corresponding
	// to this role if pass.
	UseSkill(data *types.ActionData) (bool, error)

	// Activate a passive skill corresponding to this roles.
	UsePassiveSkill(data *types.ActionData) (bool, error)
}

// Skill is used proactively.
type activeSkill struct {
	action       action.Action
	numberOfUses types.NumberOfUses
	startRound   types.Round
}

// Skill is used based on event.
type passiveSkill struct {
	action      action.Action
	afterDeath  bool
	beforeDeath bool
}

// Get roles' name.
func (r *role) GetName() string {
	return r.name
}

// Check if role has after-death passive skill.
func (r *role) HasAfterDeathSkill() bool {
	if r.passiveSkill == nil {
		return false
	}

	return r.passiveSkill.afterDeath
}

// Check if role has before-death passive skill.
func (r *role) HasBeforeDeathSkill() bool {
	if r.passiveSkill == nil {
		return false
	}

	return r.passiveSkill.beforeDeath
}

// Check condition is satisfied then perform an action corresponding
// to this role if pass.
func (r *role) UseSkill(data *types.ActionData) (bool, error) {
	if r.activeSkill == nil ||
		r.game.GetCurrentPhaseId() != r.phaseId ||
		r.game.GetCurrentRound() < r.activeSkill.startRound ||
		r.game.GetCurrentRoleId() != r.id ||
		r.activeSkill.numberOfUses == 0 {

		return false, errors.New("Unable to use skill!")
	}

	if res, err := r.activeSkill.action.Perform(data); err != nil {
		return false, err
	} else {
		r.activeSkill.numberOfUses--

		return res, nil
	}
}

// Activate a passive skill corresponding to this roles.
func (r *role) UsePassiveSkill(data *types.ActionData) (bool, error) {
	return r.passiveSkill.action.Perform(data)
}
