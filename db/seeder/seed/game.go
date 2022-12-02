package seed

import (
	"uwwolf/model"

	"gorm.io/gorm"
)

func SeedGames(client *gorm.DB) {
	client.Omit("WinningFactionID").Create([]model.Game{
		{
			ID: "1",
		},
		{
			ID: "2",
		},
	})
}
