package main

import (
	"database/sql"
	"uwwolf/app/api"
	"uwwolf/app/service"
	"uwwolf/db"
	"uwwolf/db/rdb"
	"uwwolf/util"

	_ "github.com/lib/pq"
)

func main() {
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
	rdb.ConnectRedis()

	dbi, err := sql.Open("postgres", "postgres://ww_username:ww_password@postgres/ww_db?sslmode=disable")
	if err != nil {
		panic(err)
	} else {
		db.NewStore(dbi)
	}

	roomService := service.NewRoomService(rdb.Client())
	gameService := service.NewGameService(db.DB())
	apiServer := api.NewAPIServer(roomService, gameService)
	apiServer.Run()

	// m := game.NewModerator(&game.ModeratorInit{
	// 	Scheduler:          game.NewScheduler(vars.NightPhaseID),
	// 	TurnDuration:       5 * time.Second,
	// 	DiscussionDuration: 10 * time.Second,
	// })

	// m.InitGame(&types.GameSetting{
	// 	RoleIDs:          []types.RoleID{},
	// 	NumberWerewolves: 1,
	// 	PlayerIDs: []types.PlayerID{
	// 		"1",
	// 		"2",
	// 		"3",
	// 		"4",
	// 		"5",
	// 	},
	// })
	// m.StartGame()

	// time.Sleep(3 * time.Second)

	// m.RequestPlay("1", &types.ActivateAbilityRequest{
	// 	TargetID: "2",
	// })
	// m.RequestPlay("2", &types.ActivateAbilityRequest{
	// 	TargetID: "1",
	// })
	// m.RequestPlay("3", &types.ActivateAbilityRequest{
	// 	TargetID: "4",
	// })
	// m.RequestPlay("4", &types.ActivateAbilityRequest{
	// 	TargetID: "5",
	// })
	// m.RequestPlay("5", &types.ActivateAbilityRequest{
	// 	TargetID: "3",
	// })

	// time.Sleep(8 * time.Second)
	// m.RequestPlay("1", &types.ActivateAbilityRequest{
	// 	TargetID: "2",
	// })
	// m.RequestPlay("2", &types.ActivateAbilityRequest{
	// 	TargetID: "1",
	// })
	// m.RequestPlay("3", &types.ActivateAbilityRequest{
	// 	TargetID: "4",
	// })
	// m.RequestPlay("4", &types.ActivateAbilityRequest{
	// 	TargetID: "5",
	// })
	// m.RequestPlay("5", &types.ActivateAbilityRequest{
	// 	TargetID: "3",
	// })

	// m.FinishGame()

	// time.Sleep(1 * time.Hour)
}
