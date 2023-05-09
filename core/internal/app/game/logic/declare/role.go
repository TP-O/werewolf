package declare

import (
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"
)

const (
	VillagerRoleID types.RoleID = iota + 1
	WerewolfRoleID
	HunterRoleID
	SeerRoleID
	TwoSistersRoleID
)

// Faction ID to its role IDs
var FactionID2RoleIDs = util.NewImmutableMap(map[types.FactionID][]types.RoleID{
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
var RoleSets = util.NewImmutableMap(map[types.RoleID]types.Times{
	VillagerRoleID:   UnlimitedTimes,
	WerewolfRoleID:   UnlimitedTimes,
	HunterRoleID:     Once,
	SeerRoleID:       Once,
	TwoSistersRoleID: types.Times(2),
})

// Role weights to decide main role. The highest weight is picked as the main role
var RoleWeights = util.NewImmutableMap(map[types.RoleID]uint8{
	VillagerRoleID:   0,
	TwoSistersRoleID: 1,
	SeerRoleID:       1,
	HunterRoleID:     1,
	WerewolfRoleID:   2,
})
