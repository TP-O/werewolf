package seed

import (
	"gorm.io/gorm"

	"uwwolf/app/enum"
	"uwwolf/app/model"
)

func SeedRoles(client *gorm.DB) {
	client.Create(
		[]model.Role{
			// Night
			{
				Id:        enum.TwoSistersRoleId,
				FactionId: enum.VillagerFactionId,
				PhaseId:   enum.NightPhaseId,
			},
			{
				Id:        enum.SeerRoleId,
				FactionId: enum.VillagerFactionId,
				PhaseId:   enum.NightPhaseId,
			},
			{
				Id:        enum.WerewolfRoleId,
				FactionId: enum.WerewolfFactionId,
				PhaseId:   enum.NightPhaseId,
			},

			// Day
			{
				Id:        enum.HunterRoleId,
				FactionId: enum.VillagerFactionId,
				PhaseId:   enum.DayPhaseId,
			},
			{
				Id:        enum.VillagerRoleId,
				FactionId: enum.VillagerFactionId,
				PhaseId:   enum.DayPhaseId,
			},
		})
}
