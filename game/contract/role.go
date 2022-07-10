package contract

type SkillTurn struct {
	StartFrom    uint
	NumberOfUses int
}

type Skill struct {
	Action Action
	Turn   *SkillTurn
}

type PassiveEvent struct {
	BeforeDeath         bool
	AfterDeath          bool
	BeforeBeingExecuted bool
}

type Passive struct {
	Action Action
	Event  *PassiveEvent
}

type Role interface {
	GetName() string
	GetSkill() *Skill
	GetPassive() *Passive
}
