package seed

import (
	"uwwolf/app/model"
	"uwwolf/app/types"

	"gorm.io/gorm"
)

func SeedPlayers(client *gorm.DB) {
	client.Create([]model.Player{
		{
			Id: types.PlayerId("1111111111111111111111111111"),
		},
		{
			Id: types.PlayerId("2222222222222222222222222222"),
		},
		{
			Id: types.PlayerId("3333333333333333333333333333"),
		},
		{
			Id: types.PlayerId("4444444444444444444444444444"),
		},
		{
			Id: types.PlayerId("5555555555555555555555555555"),
		},
		{
			Id: types.PlayerId("6666666666666666666666666666"),
		},
		{
			Id: types.PlayerId("7777777777777777777777777777"),
		},
		{
			Id: types.PlayerId("8888888888888888888888888888"),
		},
	})
}
