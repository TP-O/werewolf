package enum

type Limit = int8

const (
	Unlimited Limit = iota - 1
	ReachedLimit
	OneMore
)

type Priority = uint8

// Day phase
const (
	VillagerTurnPriority Priority = iota
	HunterTurnPriority
)

// Night phase
const (
	WerewolfTurnPriority Priority = iota
	SeerTurnPriority
	TwoSistersTurnPriority
)

type RoleID = uint32

const (
	VillagerRoleID RoleID = iota + 1
	WerewolfRoleID
	HunterRoleID
	SeerRoleID
	TwoSistersRoleID
)

func IsUnknownRoleID(r RoleID) bool {
	return r == 0
}

type FactionID = uint8

const (
	VillagerFactionID FactionID = iota + 1
	WerewolfFactionID
)

func IsUnknownFactionID(f FactionID) bool {
	return f == 0
}
