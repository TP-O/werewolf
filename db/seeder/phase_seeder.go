// package seeder

// import (
// 	"gorm.io/gorm"

// 	"uwwolf/app/model"
// 	"uwwolf/enum"
// )

// func seedPhases(db *gorm.DB) {
// 	insert(
// 		db,
// 		&model.Phase{},
// 		[]model.Phase{
// 			{
// 				Model: gorm.Model{
// 					ID: enum.DayPhase,
// 				},
// 				Name: "Day",
// 			},
// 			{
// 				Model: gorm.Model{
// 					ID: enum.DuskPhase,
// 				},
// 				Name: "Dusk",
// 			},
// 			{
// 				Model: gorm.Model{
// 					ID: enum.NightPhase,
// 				},
// 				Name: "Night",
// 			},
// 		})
// }
