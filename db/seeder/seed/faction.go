package seed

import (
	"gorm.io/gorm"

	"uwwolf/module/game/model"
	"uwwolf/types"
)

func SeedFactions(client *gorm.DB) {
	client.Create([]model.Faction{
		{
			ID:   types.VillagerFaction,
			Name: "Villager",
		},
		{
			ID:   types.WerewolfFaction,
			Name: "Werewolf",
		},
	})
}
