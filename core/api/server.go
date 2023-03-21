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

	gameGroup.POST("/start", func(c *gin.Context) {
		var err error
		playerID := types.PlayerID(c.GetString("playerID"))

		var payload dto.StartGameDto
		if err := c.ShouldBind(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
			return
		}

		var room dto.Room
		roomJson := queryRoom.Run(context.Background(), rdb.Client(), []string{}, playerID).String()
		if err := json.Unmarshal([]byte(roomJson), &room); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something wrong",
			})
			return
		}

		if playerID != room.OwnerID {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Not owner",
			})
			return
		}

		mod, err := game.NewModerator(&types.ModeratorInit{
			TurnDuration:       payload.TurnDuration,
			DiscussionDuration: payload.DiscussionDuration,
			GameSetting: types.GameSetting{
				RoleIDs:          payload.RoleIDs,
				RequiredRoleIDs:  payload.RequiredRoleIDs,
				NumberWerewolves: payload.NumberWerewolves,
				PlayerIDs:        room.PlayerIDs,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		gameRecord, err := db.DB().CreateGame(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to create game",
			})
			return
		}

		game.Manager().AddModerator(uint64(gameRecord.ID), mod) // nolint errcheck
		mod.StartGame()

		// Notify players

		c.JSON(http.StatusOK, gin.H{
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
