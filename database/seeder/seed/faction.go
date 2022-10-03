package seed

import (
	"gorm.io/gorm"

	"uwwolf/app/enum"
	"uwwolf/app/model"
)

func SeedFactions(client *gorm.DB) {
	client.Create([]model.Faction{
		{
			Id: enum.VillagerFactionId,
		},
		{
			Id: enum.WerewolfFactionId,
		},
	})
}
