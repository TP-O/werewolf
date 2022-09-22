package seed

import (
	"gorm.io/gorm"

	"uwwolf/module/game/model"
	"uwwolf/types"
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
