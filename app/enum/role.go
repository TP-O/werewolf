package enum

import (
	"uwwolf/app/types"
)

const (
	VillagerRoleId types.RoleId = iota + 1
	WerewolfRoleId
	HunterRoleId
	SeerRoleId
	TwoSistersRoleId
)

var WerewolfRoleIds = []types.RoleId{
	WerewolfRoleId,
}

var NonWerewolfRoleIds = []types.RoleId{
	VillagerRoleId,
	HunterRoleId,
	SeerRoleId,
	TwoSistersRoleId,
}

var RoleSets = map[types.RoleId]int{
	VillagerRoleId:   -1,
	WerewolfRoleId:   -1,
	HunterRoleId:     1,
	SeerRoleId:       1,
	TwoSistersRoleId: 2,
}
