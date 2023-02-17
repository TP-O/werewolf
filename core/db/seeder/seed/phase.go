package seed

import (
	"log"
	"uwwolf/db"
	"uwwolf/game/enum"
	"uwwolf/model"

	"github.com/gocql/gocql"
)

func SeedPhases(client *gocql.Session) error {
	batch := client.NewBatch(0)
	models := []model.Phase{
		{
			ID:   enum.NightPhaseID,
			Name: "Night",
		},
		{
			ID:   enum.DayPhaseID,
			Name: "Day",
		},
		{
			ID:   enum.DuskPhaseID,
			Name: "Dusk",
		},
	}

	for _, model := range models {
		batch.Query(
			"INSERT INTO phases (id, name) VALUES (?, ?)",
			model.ID,
			model.Name,
		)
	}

	log.Println("PhaseSeeder is completed!")

	return db.Client().ExecuteBatch(batch)
}
