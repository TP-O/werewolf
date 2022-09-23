package main

import (
	"log"
	"uwwolf/database/migration"
)

func main() {
	migration := migration.Migrations()

	if err := migration.RollbackLast(); err != nil {
		log.Fatalf("Could not rollback: %v", err)
	}

	log.Printf("Rollback has run successfully")
}
