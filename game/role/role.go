package role

import "uwwolf/contract/itf"

type role struct {
	name    string
	skill   *itf.Skill
	passive *itf.Passive
}

func (r *role) GetName() string {
	return r.name
}

func (r *role) GetSkill() *itf.Skill {
	return r.skill
}

func (r *role) GetPassive() *itf.Passive {
	return r.passive
}
