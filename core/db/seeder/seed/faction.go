package seed

// import (
// 	"log"
// 	"uwwolf/db"
// 	"uwwolf/game/enum"
// 	"uwwolf/model"

// 	"github.com/gocql/gocql"
// )

// func SeedFactions(client *gocql.Session) error {
// 	batch := client.NewBatch(0)
// 	models := []model.Faction{
// 		{
// 			ID:   enum.VillagerFactionID,
// 			Name: "Villager",
// 		},
// 		{
// 			ID:   enum.WerewolfFactionID,
// 			Name: "Werewolf",
// 		},
// 	}

// 	for _, model := range models {
// 		batch.Query(
// 			"INSERT INTO factions (id, name) VALUES (?, ?)",
// 			model.ID,
// 			model.Name,
// 		)
// 	}

// 	log.Println("FactionSeeder is completed!")

// 	return db.Client().ExecuteBatch(batch)
// }
