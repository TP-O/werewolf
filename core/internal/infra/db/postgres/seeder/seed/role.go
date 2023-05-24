package seed

// import (
// 	"log"
// 	"uwwolf/db"
// 	"uwwolf/game/enum"
// 	"uwwolf/model"

// 	"github.com/gocql/gocql"
// )

// func SeedRoles(client *gocql.Session) error {
// 	batch := client.NewBatch(0)
// 	models := []model.Role{
// 		//Night
// 		{
// 			ID:        enum.TwoSistersRoleID,
// 			FactionID: enum.VillagerFactionID,
// 			PhaseID:   enum.NightPhaseID,
// 			Name:      "Two Sisters",
// 		},
// 		{
// 			ID:        enum.SeerRoleID,
// 			FactionID: enum.VillagerFactionID,
// 			PhaseID:   enum.NightPhaseID,
// 			Name:      "Seer",
// 		},
// 		{
// 			ID:        enum.WerewolfRoleID,
// 			FactionID: enum.WerewolfFactionID,
// 			PhaseID:   enum.NightPhaseID,
// 			Name:      "Werewolf",
// 		},

// 		// Day
// 		{
// 			ID:        enum.HunterRoleID,
// 			FactionID: enum.VillagerFactionID,
// 			PhaseID:   enum.DayPhaseID,
// 			Name:      "Hunter",
// 		},
// 		{
// 			ID:        enum.VillagerRoleID,
// 			FactionID: enum.VillagerFactionID,
// 			PhaseID:   enum.DayPhaseID,
// 			Name:      "Villager",
// 		},
// 	}

// 	for _, model := range models {
// 		batch.Query(
// 			"INSERT INTO roles (id, faction_id, phase_id, name) VALUES (?, ?, ?, ?)",
// 			model.ID,
// 			model.FactionID,
// 			model.PhaseID,
// 			model.Name,
// 		)
// 	}

// 	log.Println("RoleSeeder is completed!")

// 	return db.Client().ExecuteBatch(batch)
// }
