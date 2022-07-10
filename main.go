package main

import (
	"fmt"
	"log"
	"uwwolf/game"
	"uwwolf/game/contract"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if instance, err := game.NewGameInstance(&contract.GameInstanceInit{
		GameId:             "11111111111111111111",
		Capacity:           5,
		NumberOfWerewolves: 1,
		RolePool:           []uint{3, 4, 5, 6},
	}); err == nil {
		instance.AddPlayers(
			[]string{
				"11111111111111111111",
				"11111111111111111112",
				"11111111111111111113",
				"11111111111111111114",
				"11111111111111111118",
			},
			[]uint{1, 2, 3, 4, 5},
		)

		fmt.Println(instance.Start())
	} else {
		fmt.Println(err)
	}
}
