package enum

type Limit int

const (
	Unlimited Limit = iota - 1
	ReachedLimit
	OneMore
)

type Priority int

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

type RoleID uint

const (
	VillagerRoleID RoleID = iota + 1
	WerewolfRoleID
	HunterRoleID
	SeerRoleID
	TwoSistersRoleID
)

func (r RoleID) IsUnknown() bool {
	return r == 0
}

type FactionID uint

const (
	VillagerFactionID FactionID = iota + 1
	WerewolfFactionID
)

func (f FactionID) IsUnknown() bool {
	return f == 0
}
