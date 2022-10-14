package seed

import (
	"gorm.io/gorm"

	"uwwolf/app/model"
)

func SeedGames(client *gorm.DB) {
	client.Omit("WinningFactionId").Create([]model.Game{
		{
			Id: 1,
		},
		{
			Id: 2,
		},
	})
}
