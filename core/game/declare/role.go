package declare

import (
	"uwwolf/game/types"
	"uwwolf/util/helper"
)

const (
	VillagerRoleID types.RoleID = iota + 1
	WerewolfRoleID
	HunterRoleID
	SeerRoleID
	TwoSistersRoleID
)

// Faction ID to its role IDs
var FactionID2RoleIDs = helper.NewImmutableMap(map[types.FactionID][]types.RoleID{
	WerewolfFactionID: {
		WerewolfRoleID,
	},
	VillagerFactionID: {
		VillagerRoleID,
		HunterRoleID,
		SeerRoleID,
		TwoSistersRoleID,
	},
})

// Number of maximum role's instances played in one game
var RoleSets = helper.NewImmutableMap(map[types.RoleID]types.Times{
	VillagerRoleID:   UnlimitedTimes,
	WerewolfRoleID:   UnlimitedTimes,
	HunterRoleID:     Once,
	SeerRoleID:       Once,
	TwoSistersRoleID: types.Times(2),
})

// Role weights to decide main role. The highest weight is picked as the main role
var RoleWeights = helper.NewImmutableMap(map[types.RoleID]uint8{
	VillagerRoleID:   0,
	TwoSistersRoleID: 1,
	SeerRoleID:       1,
	HunterRoleID:     1,
	WerewolfRoleID:   2,
})
