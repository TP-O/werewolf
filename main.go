package main

import (
	"time"
	"uwwolf/module/game/core"
	"uwwolf/types"
)

func main() {
	// if err := database.LoadDatabase(); err != nil {
	// 	log.Fatal("Error coneect to dabase: ", err)
	// }

	// validator.LoadValidator()
	// factory.LoadFactories()

	// if instance, err := game.NewGameInstance(&typ.GameInstanceInit{
	// 	GameId:             "11111111111111111111",
	// 	Capacity:           10,
	// 	NumberOfWerewolves: 2,
	// 	RolePool:           []int{enum.HunterRole},
	// }); err == nil {
	// 	instance.AddPlayers(
	// 		[]string{
	// 			"11111111111111111111",
	// 			"11111111111111111112",
	// 			"11111111111111111113",
	// 			"11111111111111111114",
	// 			"11111111111111111115",
	// 			"11111111111111111116",
	// 			"11111111111111111117",
	// 			"11111111111111111118",
	// 			"11111111111111111119",
	// 			"11111111111111111110",
	// 		},
	// 		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	// 	)

	// 	fmt.Println(instance.Start())
	// 	fmt.Println("======================")

	// 	instance.Do(&typ.ActionInstruction{
	// 		GameId:  "11111111111111111111",
	// 		Actor:   1,
	// 		Targets: []int{2},
	// 		Skipped: false,
	// 		Payload: []byte{},
	// 	})
	// 	// instance.Do(&typ.ActionInstruction{
	// 	// 	GameId:  "11111111111111111111",
	// 	// 	Actor:   3,
	// 	// 	Targets: []int{2},
	// 	// 	Skipped: false,
	// 	// 	Payload: []byte{},
	// 	// })

	// 	time.Sleep(3 * time.Second)

	// 	instance.Do(&typ.ActionInstruction{
	// 		GameId:  "11111111111111111111",
	// 		Actor:   3,
	// 		Targets: []int{2},
	// 		Skipped: false,
	// 		Payload: []byte{},
	// 	})

	// 	time.Sleep(5 * time.Second)
	// } else {
	// 	fmt.Println(err)
	// }

	// add := Address{
	// 	Email:  "Hehe",
	// 	Street: "!!!!",
	// }

	// data := types.GameInstanceInit{
	// 	Capacity:           10,
	// 	NumberOfWerewolves: 1,
	// }

	// var g contract.Game = core.NewGame(&types.GameData{
	// 	//
	// })

	// fmt.Println(g)

	core.NewGame(&types.GameSetting{
		NumberOfWerewolves: 2,
		RolePool: []types.RoleId{
			3, 4,
		},
		PlayerIds: []types.PlayerId{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
		TurnDuration:       2 * time.Second,
		DiscussionDuration: 5 * time.Second,
	}).Start()

	time.Sleep(10 * time.Second)
}
