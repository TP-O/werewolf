package api_test

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"uwwolf/app/data"
// 	"uwwolf/app/enum"
// 	"uwwolf/app/server/api"
// 	"uwwolf/config"
// 	mock_service "uwwolf/mock/app/service"
// 	mock_game "uwwolf/mock/game"
// 	"uwwolf/util"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// )

// func (ass ApiServiceSuite) TestStartGame() {
// 	tests := []struct {
// 		name  string
// 		setup func(ctx *gin.Context, gameSvc *mock_service.MockGameService, mod *mock_game.MockModerator)
// 		check func(res *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:  "Failure (Forget to use WaitingRoomOwner middleware)",
// 			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService, mod *mock_game.MockModerator) {},
// 			check: func(res *httptest.ResponseRecorder) {
// 				ass.Equal(http.StatusInternalServerError, res.Code)
// 				ass.Equal(
// 					"Unable to start game!",
// 					util.JsonToMap(res.Body.String())["message"],
// 				)
// 			},
// 		},
// 		{
// 			name: "Failure (Check falied)",
// 			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService, mod *mock_game.MockModerator) {
// 				room := &data.WaitingRoom{
// 					ID: "room_id",
// 				}
// 				ctx.Set(enum.WaitingRoomCtxKey, room)

// 				gameCfg := data.GameConfig{
// 					NumberWerewolves: 1,
// 				}
// 				gameSvc.EXPECT().GameConfig(room.ID).Return(gameCfg)
// 				gameSvc.EXPECT().CheckBeforeRegistration(*room, gameCfg).Return(fmt.Errorf("Check failed"))
// 			},
// 			check: func(res *httptest.ResponseRecorder) {
// 				ass.Equal(http.StatusBadRequest, res.Code)
// 				ass.Equal(
// 					"Check failed",
// 					util.JsonToMap(res.Body.String())["message"],
// 				)
// 			},
// 		},
// 		{
// 			name: "Failure (Register falied)",
// 			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService, mod *mock_game.MockModerator) {
// 				room := &data.WaitingRoom{
// 					ID: "room_id",
// 				}
// 				ctx.Set(enum.WaitingRoomCtxKey, room)

// 				gameCfg := data.GameConfig{
// 					NumberWerewolves: 1,
// 				}
// 				gameSvc.EXPECT().GameConfig(room.ID).Return(gameCfg)
// 				gameSvc.EXPECT().CheckBeforeRegistration(*room, gameCfg)
// 				gameSvc.EXPECT().RegisterGame(gameCfg, room.PlayerIDs).
// 					Return(nil, fmt.Errorf("Register failed"))
// 			},
// 			check: func(res *httptest.ResponseRecorder) {
// 				ass.Equal(http.StatusInternalServerError, res.Code)
// 				ass.Equal(
// 					"Register failed",
// 					util.JsonToMap(res.Body.String())["message"],
// 				)
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService, mod *mock_game.MockModerator) {
// 				room := &data.WaitingRoom{
// 					ID: "room_id",
// 				}
// 				ctx.Set(enum.WaitingRoomCtxKey, room)

// 				gameCfg := data.GameConfig{
// 					NumberWerewolves: 1,
// 				}
// 				gameSvc.EXPECT().GameConfig(room.ID).Return(gameCfg)
// 				gameSvc.EXPECT().CheckBeforeRegistration(*room, gameCfg)
// 				gameSvc.EXPECT().RegisterGame(gameCfg, room.PlayerIDs).Return(mod, nil)
// 				mod.EXPECT().StartGame()
// 			},
// 			check: func(res *httptest.ResponseRecorder) {
// 				ass.Equal(http.StatusOK, res.Code)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ass.Run(test.name, func() {
// 			ctrl := gomock.NewController(ass.T())
// 			defer ctrl.Finish()
// 			gameSvc := mock_service.NewMockGameService(ctrl)
// 			mod := mock_game.NewMockModerator(ctrl)

// 			res := httptest.NewRecorder()
// 			ctx, r := gin.CreateTestContext(res)

// 			test.setup(ctx, gameSvc, mod)

// 			svr := api.NewHandler(config.App{}, nil, gameSvc)
// 			r.POST("/test", func(_ *gin.Context) {
// 				svr.StartGame(ctx)
// 			})

// 			ctx.Request, _ = http.NewRequest(http.MethodPost, "/test", nil)
// 			r.ServeHTTP(res, ctx.Request)

// 			test.check(res)
// 		})
// 	}
// }
