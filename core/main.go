package main

import (
	"context"
	"fmt"
	"uwwolf/redis"
	"uwwolf/util"
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

	util.LoadConfig(".")
	// db.ConnectDB()
	c := redis.ConnectRedis()

	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	c.Set(context.Background(), "aaaa", 2, -1)
	fmt.Println(c.Get(context.Background(), "aaaa").Val())
}
