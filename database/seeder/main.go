package main

import (
	"uwwolf/database"
	"uwwolf/database/seeder/seed"
)

func main() {
	client := database.Client()

	seed.SeedPhases(client)
	seed.SeedFactions(client)
	seed.SeedRoles(client)
	seed.SeedActions(client)
	seed.SeedPlayers(client)
}
