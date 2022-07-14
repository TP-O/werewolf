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
				FactionID:   enum.VillageFaction,
				PhaseID:     enum.DayPhase,
				Name:        "Villager",
				Priority:    1,
				Score:       1,
				Quantity:    -1,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.WerewolfRole,
				},
				FactionID:   enum.WerewolfFaction,
				PhaseID:     enum.NightPhase,
				Name:        "Werewolf",
				Priority:    2,
				Score:       1,
				Quantity:    -1,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.HunterRole,
				},
				FactionID:   enum.VillageFaction,
				PhaseID:     enum.DayPhase,
				Name:        "Hunter",
				Priority:    2,
				Score:       1,
				Quantity:    1,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.SeerRole,
				},
				FactionID:   enum.VillageFaction,
				PhaseID:     enum.NightPhase,
				Name:        "Seer",
				Priority:    1,
				Score:       1,
				Quantity:    1,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.TwoSistersRole,
				},
				FactionID:   enum.VillageFaction,
				PhaseID:     enum.DayPhase,
				Name:        "Two sisters",
				Priority:    2,
				Score:       1,
				Quantity:    2,
				Image:       "image here",
				Description: "decription",
			},
			{
				Model: gorm.Model{
					ID: enum.AlphaWolfRole,
				},
				FactionID:   enum.WerewolfFaction,
				PhaseID:     enum.NightPhase,
				Name:        "Alpha Wolf",
				Priority:    3,
				Score:       1,
				Quantity:    1,
				Image:       "image here",
				Description: "decription",
			},
		})
}
