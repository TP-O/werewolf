package seeder

import (
	"gorm.io/gorm"

	"uwwolf/app/model"
	"uwwolf/enum"
)

func seedTeams(db *gorm.DB) {
	insert(
		db,
		&model.Team{},
		[]model.Team{
			{
				Model: gorm.Model{
					ID: enum.VillageFaction,
				},
				Name: "Village Faction",
			},
			{
				Model: gorm.Model{
					ID: enum.WerewolfFaction,
				},
				Name: "Werewolf Faction",
			},
			{
				Model: gorm.Model{
					ID: enum.IndependentTeam,
				},
				Name: "Independent Faction",
			},
		})

}
