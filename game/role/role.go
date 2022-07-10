package role

import "uwwolf/game/contract"

type role struct {
	name    string
	skill   *contract.Skill
	passive *contract.Passive
}

func (r *role) GetName() string {
	return r.name
}

func (r *role) GetSkill() *contract.Skill {
	return r.skill
}

func (r *role) GetPassive() *contract.Passive {
	return r.passive
}
