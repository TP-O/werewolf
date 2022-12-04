package seed

// import (
// 	"uwwolf/game/enum"
// 	"uwwolf/model"

// 	"gorm.io/gorm"
// )

// func SeedPlayers(client *gorm.DB) {
// 	client.Create([]model.Player{
// 		{
// 			ID:       enum.PlayerID("1111111111111111111111111111"),
// 			StatusID: enum.OnlineStatus,
// 		},
// 		{
// 			ID:       enum.PlayerID("2222222222222222222222222222"),
// 			StatusID: enum.OnlineStatus,
// 		},
// 		{
// 			ID:       enum.PlayerID("3333333333333333333333333333"),
// 			StatusID: enum.OnlineStatus,
// 		},
// 		{
// 			ID:       enum.PlayerID("4444444444444444444444444444"),
// 			StatusID: enum.OnlineStatus,
// 		},
// 		{
// 			ID:       enum.PlayerID("5555555555555555555555555555"),
// 			StatusID: enum.OnlineStatus,
// 		},
// 	})

// 	client.Omit("StatusID").Create([]model.Player{
// 		{
// 			ID: enum.PlayerID("6666666666666666666666666666"),
// 		},
// 		{
// 			ID: enum.PlayerID("7777777777777777777777777777"),
// 		},
// 		{
// 			ID: enum.PlayerID("8888888888888888888888888888"),
// 		},
// 	})
// }
