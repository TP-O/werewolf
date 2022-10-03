package seed

import (
	"gorm.io/gorm"

	"uwwolf/app/enum"
	"uwwolf/app/model"
)

func SeedPhases(client *gorm.DB) {
	client.Create([]model.Phase{
		{
			Id: enum.NightPhaseId,
		},
		{
			Id: enum.DayPhaseId,
		},
		{
			Id: enum.DuskPhaseId,
		},
	})
}
