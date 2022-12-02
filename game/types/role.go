package types

import "uwwolf/game/enum"

type UseRoleRequest struct {
	ActionID  enum.ActionID
	TargetIDs []enum.PlayerID
	IsSkipped bool
}

var RoleIDsByFactionID = map[enum.FactionID][]enum.RoleID{
	enum.WerewolfFactionID: {
		enum.WerewolfRoleID,
	},
	enum.VillagerFactionID: {
		enum.VillagerRoleID,
		enum.HunterRoleID,
		enum.SeerRoleID,
		enum.TwoSistersRoleID,
	},
}

var RoleIDSets = map[enum.RoleID]enum.Limit{
	enum.VillagerRoleID:   enum.Unlimited,
	enum.WerewolfRoleID:   enum.Unlimited,
	enum.HunterRoleID:     enum.OneMore,
	enum.SeerRoleID:       enum.OneMore,
	enum.TwoSistersRoleID: 2,
}

var RoleIDRanks = map[enum.RoleID]enum.Priority{
	enum.VillagerRoleID:   0,
	enum.WerewolfRoleID:   1,
	enum.TwoSistersRoleID: 2,
	enum.SeerRoleID:       3,
	enum.HunterRoleID:     4,
}
