package types

type Role int

const (
	VillagerRole Role = iota + 1
	WerewolfRole
	HunterRole
	SeerRole
	TwoSistersRole
	AlphaWolfRole
)

type SkillTurn struct {
	StartFrom    int
	NumberOfUses int
}

type PassiveEvent struct {
	BeforeDeath         bool
	AfterDeath          bool
	BeforeBeingExecuted bool
}
