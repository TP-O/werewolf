package seed

import (
	"gorm.io/gorm"

	"uwwolf/app/model"
	"uwwolf/app/types"
)

func SeedPhases(client *gorm.DB) {
	client.Create([]model.Phase{
		{
			Id:   types.NightPhase,
			Name: "Night",
		},
		{
			Id:   types.DayPhase,
			Name: "Day",
		},
		{
			Id:   types.DuskPhase,
			Name: "Dusk",
		},
	})
}
