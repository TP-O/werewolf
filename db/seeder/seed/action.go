package seed

import (
	"uwwolf/game/enum"
	"uwwolf/model"

	"gorm.io/gorm"
)

func SeedActions(client *gorm.DB) {
	client.Create([]model.Action{
		{
			ID: enum.KillActionID,
		},
		{
			ID: enum.PredictActionID,
		},
		{
			ID: enum.RecognizeActionID,
		},
		{
			ID: enum.VoteActionID,
		},
	})
}
