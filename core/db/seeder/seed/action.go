package seed

// import (
// 	"log"
// 	"uwwolf/db"
// 	"uwwolf/game/enum"
// 	"uwwolf/model"

// 	"github.com/gocql/gocql"
// )

// func SeedActions(client *gocql.Session) error {
// 	batch := client.NewBatch(0)
// 	models := []model.Action{
// 		{
// 			ID:   enum.VoteActionID,
// 			Name: "Vote",
// 		},
// 		{
// 			ID:   enum.RecognizeActionID,
// 			Name: "Recognize",
// 		},
// 		{
// 			ID:   enum.PredictActionID,
// 			Name: "Predict",
// 		},
// 		{
// 			ID:   enum.KillActionID,
// 			Name: "Kill",
// 		},
// 	}

// 	for _, model := range models {
// 		batch.Query(
// 			"INSERT INTO actions (id, name) VALUES (?, ?)",
// 			model.ID,
// 			model.Name,
// 		)
// 	}

// 	log.Println("ActionSeeder is completed!")

// 	return db.Client().ExecuteBatch(batch)
// }
