package types

type NumberOfTimes int

const (
	UnlimitedTimes NumberOfTimes = iota - 1
	OutOfTimes
	OneTimes
)

type RoleId uint

const (
	UnknownRole RoleId = iota
	VillagerRole
	WerewolfRole
	HunterRole
	SeerRole
	TwoSistersRole
	AlphaWolfRole
)

type RoleSetting struct {
	OwnerId    PlayerId
	FactionId  FactionId
	BeginRound RoundId
	Expiration NumberOfTimes
}
