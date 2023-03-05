package vars

import "uwwolf/game/types"

// Particular faction ID
const (
	VillagerFactionID types.FactionID = iota + 1
	WerewolfFactionID
)

// Particular role ID
const (
	VillagerRoleID types.RoleID = iota + 1
	WerewolfRoleID
	HunterRoleID
	SeerRoleID
	TwoSistersRoleID
)

// Particular phase ID
const (
	NightPhaseID types.PhaseID = iota + 1
	DayPhaseID
	DuskPhaseID
)

// Particular turn ID in day phase
const (
	HunterTurnID   = PreTurn
	VillagerTurnID = MidTurn
)

// Particular turn ID in night phase
const (
	SeerTurnID       = PreTurn
	TwoSistersTurnID = PreTurn
	WerewolfTurnID   = MidTurn
)

// Faction ID to its role IDs
var FactionID2RoleIDs = map[types.FactionID][]types.RoleID{
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

// Number of maximum role's instances played in one game
var RoleSets = map[types.RoleID]types.Limit{
	VillagerRoleID:   Unlimited,
	WerewolfRoleID:   Unlimited,
	HunterRoleID:     One,
	SeerRoleID:       One,
	TwoSistersRoleID: types.Limit(2),
}

// Role weights to decide main role. The highest weight is picked as the main role
var RoleWeights = map[types.RoleID]uint8{
	VillagerRoleID:   0,
	TwoSistersRoleID: 1,
	SeerRoleID:       1,
	HunterRoleID:     1,
	WerewolfRoleID:   2,
}

// Number of uses
const (
	Unlimited types.Limit = iota - 1
	ReachedLimit
	One
)

// One phase has 3 turn indexes by default
const (
	PreTurn types.TurnID = iota
	MidTurn              // Main turn
	PostTurn
)
