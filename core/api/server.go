package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"uwwolf/api/dto"
	"uwwolf/db"
	"uwwolf/db/rdb"
	"uwwolf/game"
	"uwwolf/game/types"
	"uwwolf/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var queryRoom = redis.NewScript(fmt.Sprintf(`
    local player_id = ARGV[1]
    local room_id = redis.call("GET", %v + player_id)
    if not room_id then
        return nil
    end

    return redis.call("GET", %v + room_id)
`, "pId2rId:", "room:"))

func setupRouter() *gin.Engine {
	r := gin.Default()

	gameGroup := r.Group("/game")

	gameGroup.PUT("/setting", func(ctx *gin.Context) {
		playerID := types.PlayerID(ctx.GetString("playerID"))

		var payload types.ModeratorInit
		if err := ctx.ShouldBind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
			return
		}

		var room dto.Room
		roomJson := queryRoom.Run(context.Background(), rdb.Client(), []string{}, playerID).String()
		if err := json.Unmarshal([]byte(roomJson), &room); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something wrong",
			})
			return
		}

		if playerID != room.OwnerID {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Not owner",
			})
			return
		}

		room.ModeratorInit = payload
		roomByte, err := json.Marshal(room)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something wrong",
			})
			return
		}

		rdb.Client().Set(context.Background(), fmt.Sprintf("room:%v", room.ID), string(roomByte), -1)
	})

	gameGroup.POST("/start", func(ctx *gin.Context) {
		var err error
		playerID := types.PlayerID(ctx.GetString("playerID"))

		var room dto.Room
		roomJson := queryRoom.Run(context.Background(), rdb.Client(), []string{}, playerID).String()
		if err := json.Unmarshal([]byte(roomJson), &room); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something wrong",
			})
			return
		}

		if playerID != room.OwnerID {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Not owner",
			})
			return
		}

		mod, err := game.NewModerator(&types.ModeratorInit{
			TurnDuration:       room.TurnDuration,
			DiscussionDuration: room.DiscussionDuration,
			GameSetting: types.GameSetting{
				RoleIDs:          room.RoleIDs,
				RequiredRoleIDs:  room.RequiredRoleIDs,
				NumberWerewolves: room.NumberWerewolves,
				PlayerIDs:        room.PlayerIDs,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		gameRecord, err := db.DB().CreateGame(context.Background())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to create game",
			})
			return
		}

		game.Manager().AddModerator(uint64(gameRecord.ID), mod) // nolint errcheck
		mod.StartGame()

		// Notify players

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Ok",
		})
	})

	return r
}

func StartApi() {
	r := setupRouter()

	if err := r.Run(fmt.Sprintf(":%v", util.Config().App.Port)); err != nil {
		panic(err)
	}
}
