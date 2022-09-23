package seed

import (
	"gorm.io/gorm"

	"uwwolf/app/model"
	"uwwolf/app/types"
)

func SeedFactions(client *gorm.DB) {
	client.Create([]model.Faction{
		{
			Id:   types.VillagerFaction,
			Name: "Villager",
		},
		{
			Id:   types.WerewolfFaction,
			Name: "Werewolf",
		},
	})
}
