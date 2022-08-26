package main

import (
	"uwwolf/db"
	"uwwolf/db/seeder/seed"
)

func main() {
	client := db.Client()

	seed.SeedPhases(client)
	seed.SeedFactions(client)
	seed.SeedRoles(client)
}