package seed

import (
	"gorm.io/gorm"

	"uwwolf/module/game/model"
	"uwwolf/types"
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
