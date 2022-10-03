package seed

import (
	"uwwolf/app/enum"
	"uwwolf/app/model"

	"gorm.io/gorm"
)

func SeedActions(client *gorm.DB) {
	client.Create([]model.Action{
		{
			Id: enum.ShootingActionId,
		},
		{
			Id: enum.ProphecyActionId,
		},
		{
			Id: enum.RecognitionActionId,
		},
		{
			Id: enum.VoteActionId,
		},
	})
}
