package seed

import (
	"gorm.io/gorm"

	"uwwolf/app/model"
	"uwwolf/app/types"
)

func SeedRoles(client *gorm.DB) {
	client.Create(
		[]model.Role{
			// Night
			{
				Id:         types.TwoSistersRole,
				FactionId:  types.VillagerFaction,
				PhaseId:    types.NightPhase,
				IsDefault:  true,
				Priority:   1,
				Weight:     1,
				Set:        2,
				BeginRound: 1,
				Expiration: types.OneTimes,
			},
			{
				Id:         types.SeerRole,
				FactionId:  types.VillagerFaction,
				PhaseId:    types.NightPhase,
				IsDefault:  true,
				Priority:   2,
				Weight:     1,
				Set:        1,
				BeginRound: 1,
				Expiration: types.UnlimitedTimes,
			},
			{
				Id:         types.WerewolfRole,
				FactionId:  types.WerewolfFaction,
				PhaseId:    types.NightPhase,
				IsDefault:  true,
				Priority:   3,
				Weight:     1,
				Set:        -1,
				BeginRound: 1,
				Expiration: types.UnlimitedTimes,
			},

			// Day
			{
				Id:         types.HunterRole,
				FactionId:  types.VillagerFaction,
				PhaseId:    types.DayPhase,
				IsDefault:  false,
				Priority:   0,
				Weight:     1,
				Set:        1,
				BeginRound: 1,
				Expiration: types.OneTimes,
			},
			{
				Id:         types.VillagerRole,
				FactionId:  types.VillagerFaction,
				PhaseId:    types.DayPhase,
				IsDefault:  true,
				Priority:   1,
				Weight:     1,
				Set:        -1,
				BeginRound: 1,
				Expiration: types.UnlimitedTimes,
			},
		})
}
