package seed

import (
	"log"
	"uwwolf/db"
	"uwwolf/game/enum"
	"uwwolf/model"

	"github.com/gocql/gocql"
)

func SeedPlayerStatus(client *gocql.Session) error {
	batch := client.NewBatch(0)
	models := []model.PlayerStatus{
		{
			ID:   enum.OfflineStatus,
			Name: "Offline",
		},
		{
			ID:   enum.OnlineStatus,
			Name: "Online",
		},
		{
			ID:   enum.BusyStatus,
			Name: "Busy",
		},
		{
			ID:   enum.InGameStatus,
			Name: "In game",
		},
	}

	for _, model := range models {
		batch.Query(
			"INSERT INTO player_status (id, name) VALUES (?, ?)",
			model.ID,
			model.Name,
		)
	}

	log.Println("PlayerStatusSeeder is completed!")

	return db.Client().ExecuteBatch(batch)
}
