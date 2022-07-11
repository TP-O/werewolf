package seeder

import (
	"gorm.io/gorm"

	"uwwolf/app/model"
	"uwwolf/enum"
)

func seedRoles(db *gorm.DB) {
	insert(
		db,
		&model.Role{},
		[]model.Role{
			{
				Model: gorm.Model{
					ID: enum.VillagerRole,
				},
				TeamID:      enum.VillageFaction,
				Name:        "Villager",
				Score:       1,
				Quantity:    -1,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.WerewolfRole,
				},
				TeamID:      enum.WerewolfFaction,
				Name:        "Werewolf",
				Score:       1,
				Quantity:    -1,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.HunterRole,
				},
				TeamID:      enum.VillageFaction,
				Name:        "Hunter",
				Score:       1,
				Quantity:    1,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.SeerRole,
				},
				TeamID:      enum.VillageFaction,
				Name:        "Seer",
				Score:       1,
				Quantity:    1,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.TwoSistersRole,
				},
				TeamID:      enum.VillageFaction,
				Name:        "Two sisters",
				Score:       1,
				Quantity:    2,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.AlphaWerewolfRole,
				},
				TeamID:      enum.WerewolfFaction,
				Name:        "Alpha Werewolf",
				Score:       1,
				Quantity:    1,
				Image:       "image here",
				Description: "decription",
			},
		})
}
