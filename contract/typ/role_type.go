package typ

type SkillTurn struct {
	StartFrom    uint
	NumberOfUses int
}

type PassiveEvent struct {
	BeforeDeath         bool
	AfterDeath          bool
	BeforeBeingExecuted bool
}
