package constants

import (
	"uwwolf/internal/app/game/logic/types"
)

const (
	VillagerRoleId types.RoleId = iota + 1
	WerewolfRoleId
	HunterRoleId
	SeerRoleId
	TwoSistersRoleId
)

// Faction Id to its role Ids
var FactionId2RoleIds = NewImmutableMap(map[types.FactionId][]types.RoleId{
	WerewolfFactionId: {
		WerewolfRoleId,
	},
	VillagerFactionId: {
		VillagerRoleId,
		HunterRoleId,
		SeerRoleId,
		TwoSistersRoleId,
	},
})

// Number of maximum role's instances played in one game
var RoleSets = NewImmutableMap(map[types.RoleId]types.Times{
	VillagerRoleId:   UnlimitedTimes,
	WerewolfRoleId:   UnlimitedTimes,
	HunterRoleId:     Once,
	SeerRoleId:       Once,
	TwoSistersRoleId: Twice,
})

// Role weights to decide main role. The highest weight is picked as the main role
var RoleWeights = NewImmutableMap(map[types.RoleId]uint8{
	VillagerRoleId:   0,
	TwoSistersRoleId: 1,
	SeerRoleId:       1,
	HunterRoleId:     1,
	WerewolfRoleId:   2,
})
