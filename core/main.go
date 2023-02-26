package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// grpc.Start()

	// fmt.Println(validator.ValidateStruct(types.GameSetting{
	// 	TurnDuration:       50,
	// 	DiscussionDuration: 90,
	// 	RoleIDs:            []enum.RoleID{1, 2},
	// 	NumberWerewolves:   1,
	// 	PlayerIDs: []enum.PlayerID{
	// 		"11111111111111111111",
	// 		"22222222222222222222",
	// 		"33333333333333333333",
	// 		"44444444444444444444",
	// 		"55555555555555555555",
	// 	},
	// }))

	b := BB{
		AA: "cc",
	}
	jsonStr, _ := json.Marshal(b)
	fmt.Println(string(jsonStr))
}

type AA interface {
	hello()
}

type BB struct {
	AA string `json:"aa"`
}

func (b BB) hello() {
	//
}
