package seed

import (
	"uwwolf/app/enum"
	"uwwolf/app/model"

	"gorm.io/gorm"
)

func SeedStatus(client *gorm.DB) {
	client.Create([]model.Status{
		{
			Id: enum.OnlineStatus,
		},
		{
			Id: enum.BusyStatus,
		},
		{
			Id: enum.InGameStatus,
		},
	})
}
