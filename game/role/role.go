package role

import (
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
)

type role struct {
	name    string
	skill   *itf.Skill
	passive *itf.Passive
	game    itf.IGame
}

func (r *role) GetName() string {
	return r.name
}

func (r *role) UseSkill(instruction *typ.ActionInstruction) bool {
	return r.skill.Action.Perform(r.game, instruction)
}

func (r *role) ActivatePassive(instruction *typ.ActionInstruction) bool {
	return r.passive.Action.Perform(r.game, instruction)
}
