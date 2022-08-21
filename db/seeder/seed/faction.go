package seed

import (
	"gorm.io/gorm"

	"uwwolf/module/game/model"
	"uwwolf/types"
)

func SeedFactions(client *gorm.DB) {
	client.Create([]model.Faction{
		{
			ID:   types.VillageFaction,
			Name: "Village",
		},
		{
			ID:   types.WerewolfFaction,
			Name: "Werewolf",
		},
		{
			ID:   types.IndependentFaction,
			Name: "Independent",
		},
	})
}
