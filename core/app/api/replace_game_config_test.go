package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"uwwolf/app/api"
	"uwwolf/app/data"
	"uwwolf/app/dto"
	"uwwolf/app/enum"
	"uwwolf/config"
	mock_service "uwwolf/mock/app/service"
	"uwwolf/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func (ass ApiServiceSuite) TestReplaceGameConfig() {
	tests := []struct {
		name    string
		payload string
		setup   func(ctx *gin.Context, gameSvc *mock_service.MockGameService)
		check   func(res *httptest.ResponseRecorder)
	}{
		{
			name: "Failure (Invalid request - Empty RoleIDs)",
			payload: `{
                "role_ids": [],
                "required_role_ids": [3],
                "number_werewolves": 1,
                "turn_duration": 50,
                "discussion_duration": 120
            }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"role_ids must contain at least 1 item",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["role_ids"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Invalid RoleIDs)",
			payload: `{
                "role_ids": [1],
                "required_role_ids": [3],
                "number_werewolves": 1,
                "turn_duration": 50,
                "discussion_duration": 120
            }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"role_ids[0] must be greater than 2",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["role_ids[0]"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Duplicate RolelIDs)",
			payload: `{
                "role_ids": [3, 3],
                "required_role_ids": [3],
                "number_werewolves": 1,
                "turn_duration": 50,
                "discussion_duration": 120
            }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"role_ids must contain unique values",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["role_ids"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Invalid RequiredRoleIDs)",
			payload: `{
                "role_ids": [3, 4],
                "required_role_ids": [5],
                "number_werewolves": 1,
                "turn_duration": 50,
                "discussion_duration": 120
            }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"required_role_ids must be in role_ids",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["required_role_ids"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Empty NumberWerewolves)",
			payload: `{
                "role_ids": [3, 4],
                "required_role_ids": [3],
                "number_werewolves": 0,
                "turn_duration": 50,
                "discussion_duration": 120
            }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"number_werewolves is a required field",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["number_werewolves"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Too large NumberWerewolves)",
			payload: `{
		        "role_ids": [3, 4],
		        "required_role_ids": [3],
		        "number_werewolves": 9,
		        "turn_duration": 50,
		        "discussion_duration": 120
		    }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"number_werewolves must be less than 8",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["number_werewolves"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Too small TurnDuration)",
			payload: `{
		        "role_ids": [3, 4],
		        "required_role_ids": [3],
		        "number_werewolves": 1,
		        "turn_duration": 10,
		        "discussion_duration": 120
		    }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"turn_duration must be from 20 to 60 seconds",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["turn_duration"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Too large TurnDuration)",
			payload: `{
		        "role_ids": [3, 4],
		        "required_role_ids": [3],
		        "number_werewolves": 1,
		        "turn_duration": 90,
		        "discussion_duration": 120
		    }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"turn_duration must be from 20 to 60 seconds",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["turn_duration"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Too small DiscussionDuration)",
			payload: `{
		        "role_ids": [3, 4],
		        "required_role_ids": [3],
		        "number_werewolves": 1,
		        "turn_duration": 50,
		        "discussion_duration": 10
		    }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"discussion_duration must be from 40 to 360 seconds",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["discussion_duration"],
				)
			},
		},
		{
			name: "Failure (Invalid request - Too large DiscussionDuration)",
			payload: `{
		        "role_ids": [3, 4],
		        "required_role_ids": [3],
		        "number_werewolves": 1,
		        "turn_duration": 50,
		        "discussion_duration": 500
		    }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusBadRequest, res.Code)
				ass.Equal(
					"discussion_duration must be from 40 to 360 seconds",
					util.JsonToMap(res.Body.String())["errors"].(map[string]any)["discussion_duration"],
				)
			},
		},
		{
			name: "Failure (Forget to use WaitingRoomOwner middleware)",
			payload: `{
		        "role_ids": [3, 4],
		        "required_role_ids": [3],
		        "number_werewolves": 1,
		        "turn_duration": 50,
		        "discussion_duration": 90
		    }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusInternalServerError, res.Code)
				ass.Equal(
					"Unable to update game config!",
					util.JsonToMap(res.Body.String())["message"],
				)
			},
		},
		{
			name: "Failure (Updated falied)",
			payload: `{
		        "role_ids": [3, 4],
		        "required_role_ids": [3],
		        "number_werewolves": 1,
		        "turn_duration": 50,
		        "discussion_duration": 90
		    }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {
				room := &data.WaitingRoom{
					ID: "room_id",
				}
				ctx.Set(enum.WaitingRoomCtxKey, room)

				var gameCfg dto.ReplaceGameConfigDto
				json.Unmarshal([]byte(`{
                    "role_ids": [3, 4],
                    "required_role_ids": [3],
                    "number_werewolves": 1,
                    "turn_duration": 50,
                    "discussion_duration": 90
                }`), &gameCfg)

				gameSvc.EXPECT().UpdateGameConfig(room.ID, gameCfg).Return(fmt.Errorf("Update failed"))
			},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusInternalServerError, res.Code)
				ass.Equal(
					"Something went wrong!",
					util.JsonToMap(res.Body.String())["message"],
				)
			},
		},
		{
			name: "Failure (Updated falied)",
			payload: `{
		        "role_ids": [3, 4],
		        "required_role_ids": [3],
		        "number_werewolves": 1,
		        "turn_duration": 50,
		        "discussion_duration": 90
		    }`,
			setup: func(ctx *gin.Context, gameSvc *mock_service.MockGameService) {
				room := &data.WaitingRoom{
					ID: "room_id",
				}
				ctx.Set(enum.WaitingRoomCtxKey, room)

				var gameCfg dto.ReplaceGameConfigDto
				json.Unmarshal([]byte(`{
                    "role_ids": [3, 4],
                    "required_role_ids": [3],
                    "number_werewolves": 1,
                    "turn_duration": 50,
                    "discussion_duration": 90
                }`), &gameCfg)

				gameSvc.EXPECT().UpdateGameConfig(room.ID, gameCfg)
			},
			check: func(res *httptest.ResponseRecorder) {
				ass.Equal(http.StatusOK, res.Code)
			},
		},
	}

	for _, test := range tests {
		ass.Run(test.name, func() {
			ctrl := gomock.NewController(ass.T())
			defer ctrl.Finish()
			gameSvc := mock_service.NewMockGameService(ctrl)

			res := httptest.NewRecorder()
			ctx, r := gin.CreateTestContext(res)

			test.setup(ctx, gameSvc)

			svr := api.NewServer(config.App{}, nil, gameSvc)
			r.POST("/test", func(_ *gin.Context) {
				svr.ReplaceGameConfig(ctx)
			})

			ctx.Request, _ = http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(test.payload))
			r.ServeHTTP(res, ctx.Request)

			test.check(res)
		})
	}
}
