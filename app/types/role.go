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
)

type RoleSetting struct {
	Id         RoleId
	OwnerId    PlayerId
	FactionId  FactionId
	PhaseId    PhaseId
	BeginRound RoundId
	Expiration NumberOfTimes
}
