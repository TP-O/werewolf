package config

import "uwwolf/app/game/types"

const (
	VillagerRoleID types.RoleID = iota + 1
	WerewolfRoleID
	HunterRoleID
	SeerRoleID
	TwoSistersRoleID
)

const (
	VillagerFactionID types.FactionID = iota + 1
	WerewolfFactionID
)

var RolesByFaction = map[types.FactionID][]types.RoleID{
	WerewolfFactionID: {
		WerewolfRoleID,
	},
	VillagerFactionID: {
		VillagerRoleID,
		HunterRoleID,
		SeerRoleID,
		TwoSistersRoleID,
	},
}

const (
	Unlimited types.Limit = iota - 1
	ReachedLimit
	OneMore
)

var RoleSets = map[types.RoleID]types.Limit{
	VillagerRoleID:   Unlimited,
	WerewolfRoleID:   Unlimited,
	HunterRoleID:     OneMore,
	SeerRoleID:       OneMore,
	TwoSistersRoleID: 2,
}

// Day
const (
	VillagerTurnPriority types.Priority = iota
	HunterTurnPriority
)

// Night
const (
	WerewolfTurnPriority types.Priority = iota
	SeerTurnPriority
	TwoSistersTurnPriority
)

var RolePriorities = map[types.RoleID]types.Priority{
	VillagerRoleID:   0,
	WerewolfRoleID:   1,
	TwoSistersRoleID: 2,
	SeerRoleID:       3,
	HunterRoleID:     4,
}
