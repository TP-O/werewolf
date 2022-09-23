package main

import (
	"log"
	"uwwolf/database/migration"
)

func main() {
	migration := migration.Migrations()

	if err := migration.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Printf("Migration has run successfully")
}
