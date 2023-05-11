package constants

import (
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"
)

const (
	VillagerRoleId types.RoleId = iota + 1
	WerewolfRoleId
	HunterRoleId
	SeerRoleId
	TwoSistersRoleId
)

// Faction Id to its role Ids
var FactionId2RoleIds = util.NewImmutableMap(map[types.FactionId][]types.RoleId{
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
var RoleSets = util.NewImmutableMap(map[types.RoleId]types.Times{
	VillagerRoleId:   UnlimitedTimes,
	WerewolfRoleId:   UnlimitedTimes,
	HunterRoleId:     Once,
	SeerRoleId:       Once,
	TwoSistersRoleId: types.Times(2),
})

// Role weights to decIde main role. The highest weight is picked as the main role
var RoleWeights = util.NewImmutableMap(map[types.RoleId]uint8{
	VillagerRoleId:   0,
	TwoSistersRoleId: 1,
	SeerRoleId:       1,
	HunterRoleId:     1,
	WerewolfRoleId:   2,
})
