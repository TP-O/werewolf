package seed

import (
	"gorm.io/gorm"

	"uwwolf/module/game/model"
	"uwwolf/types"
)

func SeedRoles(client *gorm.DB) {
	client.Create(
		[]model.Role{
			// Night
			{
				ID:         types.TwoSistersRole,
				FactionID:  types.VillagerFaction,
				PhaseID:    types.NightPhase,
				Name:       "Two sisters",
				Priority:   1,
				Weight:     1,
				Set:        2,
				BeginRound: 1,
				Expiration: types.OneTimes,
			},
			{
				ID:         types.SeerRole,
				FactionID:  types.VillagerFaction,
				PhaseID:    types.NightPhase,
				Name:       "Seer",
				Priority:   2,
				Weight:     1,
				Set:        1,
				BeginRound: 1,
				Expiration: types.UnlimitedTimes,
			},
			{
				ID:         types.WerewolfRole,
				FactionID:  types.WerewolfFaction,
				PhaseID:    types.NightPhase,
				Name:       "Werewolf",
				Priority:   3,
				Weight:     1,
				Set:        -1,
				BeginRound: 1,
				Expiration: types.UnlimitedTimes,
			},

			// Day
			{
				ID:         types.HunterRole,
				FactionID:  types.VillagerFaction,
				PhaseID:    types.DayPhase,
				Name:       "Hunter",
				Priority:   0,
				Weight:     1,
				Set:        1,
				BeginRound: 1,
				Expiration: types.OneTimes,
			},
			{
				ID:         types.VillagerRole,
				FactionID:  types.VillagerFaction,
				PhaseID:    types.DayPhase,
				Name:       "Villager",
				Priority:   1,
				Weight:     1,
				Set:        -1,
				BeginRound: 1,
				Expiration: types.UnlimitedTimes,
			},
		})
}
