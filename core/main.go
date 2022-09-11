package main

import (
	"fmt"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	"uwwolf/grpc"
	"uwwolf/validator"
)

func main() {
	grpc.Start()

	fmt.Println(validator.ValidateStruct(types.GameSetting{
		TurnDuration:       50,
		DiscussionDuration: 90,
		RoleIDs:            []enum.RoleID{1, 2},
		NumberWerewolves:   1,
		PlayerIDs: []enum.PlayerID{
			"11111111111111111111",
			"22222222222222222222",
			"33333333333333333333",
			"44444444444444444444",
			"55555555555555555555",
		},
	}))
}