package main

import "fmt"

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

	p := make(map[string][]map[string]*VoteRecord)
	p["aa"] = make([]map[string]*VoteRecord, 1)
	p["aa"][0] = map[string]*VoteRecord{
		"b": {
			Votes: 999,
		},
	}

	for _, aa := range p["aa"][0] {
		aa.Votes = 0
	}

	fmt.Println(p["aa"][0]["c"])
}

type IAA interface {
	hello()
}

type AA struct {
	//
}

// func (a AA) hello() {}

var _ IAA = (*AA)(nil)
