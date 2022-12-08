package seed

import (
	"log"
	"uwwolf/db"
	"uwwolf/model"

	"github.com/gocql/gocql"
)

func SeedPlayers(client *gocql.Session) error {
	batch := client.NewBatch(0)
	models := []model.Player{
		{
			ID:       "11111111111111111111",
			StatusID: 0,
		},
		{
			ID:       "22222222222222222222",
			StatusID: 0,
		},
		{
			ID:       "33333333333333333333",
			StatusID: 0,
		},
		{
			ID:       "44444444444444444444",
			StatusID: 0,
		},
		{
			ID:       "55555555555555555555",
			StatusID: 0,
		},
	}

	for _, model := range models {
		batch.Query(
			"INSERT INTO players (id, status_id) VALUES (?, ?)",
			model.ID,
			model.StatusID,
		)
	}

	log.Println("PlayerSeeder is completed!")

	return db.Client().ExecuteBatch(batch)
}
