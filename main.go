package main

import (
	"fmt"
	"log"

	"uwwolf/config"
	"uwwolf/contract/typ"
	"uwwolf/database"
	"uwwolf/game"
	"uwwolf/validator"
)

func main() {
	if err := config.LoadConfigs(); err != nil {
		log.Fatal("Error loading config: ", err)
	}

	if err := database.LoadDatabase(); err != nil {
		log.Fatal("Error coneect to dabase: ", err)
	}

	validator.LoadValidator()

	if instance, err := game.NewGameInstance(&typ.GameInstanceInit{
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
