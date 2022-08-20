// package seeder

// import (
// 	"gorm.io/gorm"

// 	"uwwolf/app/model"
// 	"uwwolf/enum"
// )

// func seedFactions(db *gorm.DB) {
// 	insert(
// 		db,
// 		&model.Faction{},
// 		[]model.Faction{
// 			{
// 				Model: gorm.Model{
// 					ID: enum.VillageFaction,
// 				},
// 				Name: "Village Faction",
// 			},
// 			{
// 				Model: gorm.Model{
// 					ID: enum.WerewolfFaction,
// 				},
// 				Name: "Werewolf Faction",
// 			},
// 			{
// 				Model: gorm.Model{
// 					ID: enum.IndependentTeam,
// 				},
// 				Name: "Independent Faction",
// 			},
// 		})

// }
