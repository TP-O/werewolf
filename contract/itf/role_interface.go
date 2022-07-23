package itf

import "uwwolf/contract/typ"

type Skill struct {
	Action IAction
	Turn   *typ.SkillTurn
}

type Passive struct {
	Action IAction
	Event  *typ.PassiveEvent
}

type IRole interface {
	GetName() string
	UseSkill(instruction *typ.ActionInstruction) bool
	ActivatePassive(instruction *typ.ActionInstruction) bool
}
