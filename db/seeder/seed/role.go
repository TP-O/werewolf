package seed

import (
	"gorm.io/gorm"

	"uwwolf/game/enum"
	"uwwolf/model"
)

func SeedRoles(client *gorm.DB) {
	client.Create(
		[]model.Role{
			// Night
			{
				ID:        enum.TwoSistersRoleID,
				FactionID: enum.VillagerFactionID,
				PhaseID:   enum.NightPhaseID,
			},
			{
				ID:        enum.SeerRoleID,
				FactionID: enum.VillagerFactionID,
				PhaseID:   enum.NightPhaseID,
			},
			{
				ID:        enum.WerewolfRoleID,
				FactionID: enum.WerewolfFactionID,
				PhaseID:   enum.NightPhaseID,
			},

			// Day
			{
				ID:        enum.HunterRoleID,
				FactionID: enum.VillagerFactionID,
				PhaseID:   enum.DayPhaseID,
			},
			{
				ID:        enum.VillagerRoleID,
				FactionID: enum.VillagerFactionID,
				PhaseID:   enum.DayPhaseID,
			},
		})
}
