package seed

import (
	"gorm.io/gorm"

	"uwwolf/module/game/model"
	"uwwolf/types"
)

func SeedPhases(client *gorm.DB) {
	client.Create([]model.Phase{
		{
			ID:   types.NightPhase,
			Name: "Night",
		},
		{
			ID:   types.DayPhase,
			Name: "Day",
		},
		{
			ID:   types.DuskPhase,
			Name: "Dusk",
		},
	})
}
