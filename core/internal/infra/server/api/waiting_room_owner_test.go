package api_test

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"uwwolf/app/data"
// 	"uwwolf/app/enum"
// 	"uwwolf/app/server/api"
// 	"uwwolf/config"
// 	mock_service "uwwolf/mock/app/service"
// 	"uwwolf/util"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// )

// func (ass ApiServiceSuite) TestWaitingRoomOwner() {
// 	tests := []struct {
// 		name  string
// 		setup func(roomSvc *mock_service.MockRoomService)
// 		check func(ctx *gin.Context, res *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "Failure (Non-existent room)",
// 			setup: func(roomSvc *mock_service.MockRoomService) {
// 				roomSvc.EXPECT().PlayerWaitingRoom(ass.playerID1).Return(data.WaitingRoom{}, false)
// 			},
// 			check: func(ctx *gin.Context, res *httptest.ResponseRecorder) {
// 				ass.Nil(ctx.Get(enum.WaitingRoomCtxKey))
// 				ass.Equal(http.StatusBadRequest, res.Code)
// 				ass.Equal("You're not in any room!", util.JsonToMap(res.Body.String())["message"])
// 			},
// 		},
// 		{
// 			name: "Failure (Not room owner)",
// 			setup: func(roomSvc *mock_service.MockRoomService) {
// 				roomSvc.EXPECT().PlayerWaitingRoom(ass.playerID1).Return(data.WaitingRoom{
// 					OwnerID: ass.playerID2,
// 				}, true)
// 			},
// 			check: func(ctx *gin.Context, res *httptest.ResponseRecorder) {
// 				ass.Nil(ctx.Get(enum.WaitingRoomCtxKey))
// 				ass.Equal(http.StatusForbidden, res.Code)
// 				ass.Equal("Only the room owner can start the game!", util.JsonToMap(res.Body.String())["message"])
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			setup: func(roomSvc *mock_service.MockRoomService) {
// 				roomSvc.EXPECT().PlayerWaitingRoom(ass.playerID1).Return(data.WaitingRoom{
// 					OwnerID: ass.playerID1,
// 				}, true)
// 			},
// 			check: func(ctx *gin.Context, res *httptest.ResponseRecorder) {
// 				ass.NotNil(ctx.Get(enum.WaitingRoomCtxKey))
// 				ass.Equal(http.StatusOK, res.Code)
// 				ass.Nil(util.JsonToMap(res.Body.String())["message"])
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ass.Run(test.name, func() {
// 			ctrl := gomock.NewController(ass.T())
// 			defer ctrl.Finish()
// 			roomSvc := mock_service.NewMockRoomService(ctrl)

// 			test.setup(roomSvc)

// 			res := httptest.NewRecorder()
// 			ctx, r := gin.CreateTestContext(res)

// 			svr := api.NewHandler(config.App{}, roomSvc, nil)
// 			r.POST("/test", func(_ *gin.Context) {
// 				ctx.Set(enum.PlayerIDCtxKey, string(ass.playerID1))
// 				svr.WaitingRoomOwner(ctx)
// 			})

// 			ctx.Request, _ = http.NewRequest(http.MethodPost, "/test", nil)
// 			r.ServeHTTP(res, ctx.Request)

// 			test.check(ctx, res)
// 		})
// 	}
// }
