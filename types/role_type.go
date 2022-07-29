package types

type Role int

const (
	UnkownRole Role = iota
	VillagerRole
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
