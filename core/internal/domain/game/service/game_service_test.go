package service_test

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"reflect"
// 	"testing"
// 	"uwwolf/app/data"
// 	"uwwolf/app/dto"
// 	"uwwolf/app/enum"
// 	"uwwolf/app/service"
// 	"uwwolf/config"
// 	"uwwolf/db/postgres"
// 	"uwwolf/internal/app/game/logic/types"
// 	"uwwolf/game/vars"
// 	mock_storage "uwwolf/mock/db"
// 	mock_game "uwwolf/mock/game"

// 	"github.com/go-redis/redismock/v9"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type GameServiceSuite struct {
// 	suite.Suite
// 	roomID          string
// 	gameID          int64
// 	playerID        types.PlayerID
// 	joinedPlayerIDs []types.PlayerID
// }

// func (gss *GameServiceSuite) SetupSuite() {
// 	gss.roomID = "room_id"
// 	gss.gameID = 1
// 	gss.playerID = types.PlayerID("1")
// 	gss.joinedPlayerIDs = []types.PlayerID{
// 		types.PlayerID("1"),
// 		types.PlayerID("2"),
// 		types.PlayerID("3"),
// 	}
// }

// func TestGameServiceSuite(t *testing.T) {
// 	suite.Run(t, new(GameServiceSuite))
// }

// func (gss GameServiceSuite) TestNewGameService() {
// 	ctrl := gomock.NewController(gss.T())
// 	defer ctrl.Finish()
// 	pdb := mock_storage.NewMockStore(ctrl)
// 	rdb, _ := redismock.NewClusterMock()
// 	manager := mock_game.NewMockManager(ctrl)

// 	svc := service.NewGameService(config.Game{}, rdb, pdb, manager)

// 	gss.Require().NotNil(svc)
// 	gss.Require().False(reflect.ValueOf(svc).Elem().FieldByName("rdb").IsNil())
// 	gss.Require().False(reflect.ValueOf(svc).Elem().FieldByName("pdb").IsNil())
// 	gss.Require().False(reflect.ValueOf(svc).Elem().FieldByName("gameManager").IsNil())
// }

// func (gss GameServiceSuite) TestGameConfig() {
// 	tests := []struct {
// 		name           string
// 		expectedResult data.GameConfig
// 		setup          func(rdb redismock.ClusterClientMock)
// 	}{
// 		{
// 			name: "Ok (Return default config)",
// 			expectedResult: data.GameConfig{
// 				RoleIDs:            []types.RoleID{vars.SeerRoleID},
// 				NumberWerewolves:   1,
// 				TurnDuration:       20,
// 				DiscussionDuration: 90,
// 			},
// 			setup: func(rdb redismock.ClusterClientMock) {
// 				rdb.ExpectGet(enum.RoomID2GameConfigRNs + gss.roomID).
// 					SetVal("Invalid json")
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			expectedResult: data.GameConfig{
// 				RoleIDs:            []types.RoleID{vars.SeerRoleID},
// 				NumberWerewolves:   1,
// 				TurnDuration:       20,
// 				DiscussionDuration: 90,
// 			},
// 			setup: func(rdb redismock.ClusterClientMock) {
// 				encodedExpectedResult, _ := json.Marshal(data.GameConfig{
// 					RoleIDs:            []types.RoleID{vars.SeerRoleID},
// 					NumberWerewolves:   1,
// 					TurnDuration:       20,
// 					DiscussionDuration: 90,
// 				})

// 				rdb.ExpectGet(enum.RoomID2GameConfigRNs + gss.roomID).
// 					SetVal(string(encodedExpectedResult))
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		gss.Run(test.name, func() {
// 			ctrl := gomock.NewController(gss.T())
// 			defer ctrl.Finish()
// 			pdb := mock_storage.NewMockStore(ctrl)
// 			rdb, rmock := redismock.NewClusterMock()
// 			manager := mock_game.NewMockManager(ctrl)

// 			test.setup(rmock)

// 			svc := service.NewGameService(config.Game{}, rdb, pdb, manager)
// 			config := svc.GameConfig(gss.roomID)

// 			gss.Require().Equal(test.expectedResult, config)
// 		})
// 	}
// }

// func (gss GameServiceSuite) TestUpdateGameConfig() {
// 	tests := []struct {
// 		name        string
// 		config      dto.ReplaceGameConfigDto
// 		expectedErr error
// 		setup       func(rdb redismock.ClusterClientMock)
// 	}{
// 		{
// 			name:        "Failure (Store failed)",
// 			config:      dto.ReplaceGameConfigDto{},
// 			expectedErr: fmt.Errorf("Set failed"),
// 			setup: func(rdb redismock.ClusterClientMock) {
// 				encodedConfig, _ := json.Marshal(dto.ReplaceGameConfigDto{})

// 				rdb.ExpectSet(enum.RoomID2GameConfigRNs+gss.roomID, string(encodedConfig), -1).
// 					SetErr(fmt.Errorf("Set failed"))
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			config: dto.ReplaceGameConfigDto{
// 				NumberWerewolves:   1,
// 				TurnDuration:       30,
// 				DiscussionDuration: 120,
// 				RoleIDs:            []types.RoleID{vars.HunterRoleID},
// 			},
// 			expectedErr: nil,
// 			setup: func(rdb redismock.ClusterClientMock) {
// 				encodedConfig, _ := json.Marshal(dto.ReplaceGameConfigDto{
// 					NumberWerewolves:   1,
// 					TurnDuration:       30,
// 					DiscussionDuration: 120,
// 					RoleIDs:            []types.RoleID{vars.HunterRoleID},
// 				})

// 				rdb.ExpectSet(enum.RoomID2GameConfigRNs+gss.roomID, string(encodedConfig), -1).
// 					SetVal("Ok")
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		gss.Run(test.name, func() {
// 			ctrl := gomock.NewController(gss.T())
// 			defer ctrl.Finish()
// 			pdb := mock_storage.NewMockStore(ctrl)
// 			rdb, rmock := redismock.NewClusterMock()
// 			manager := mock_game.NewMockManager(ctrl)

// 			test.setup(rmock)

// 			svc := service.NewGameService(config.Game{}, rdb, pdb, manager)
// 			err := svc.UpdateGameConfig(gss.roomID, test.config)

// 			gss.Require().Equal(test.expectedErr, err)
// 		})
// 	}
// }

// func (gss GameServiceSuite) TestCheckBeforeRegistration() {
// 	tests := []struct {
// 		name        string
// 		room        data.WaitingRoom
// 		cfg         data.GameConfig
// 		expectedErr error
// 	}{
// 		{
// 			name: "Not enough players",
// 			room: data.WaitingRoom{
// 				PlayerIDs: []types.PlayerID{"1"},
// 			},
// 			cfg:         data.GameConfig{},
// 			expectedErr: fmt.Errorf("Invite more players to play!"),
// 		},
// 		{
// 			name: "Too many players",
// 			room: data.WaitingRoom{
// 				PlayerIDs: []types.PlayerID{
// 					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11",
// 				},
// 			},
// 			cfg:         data.GameConfig{},
// 			expectedErr: fmt.Errorf("Too many players!"),
// 		},
// 		{
// 			name: "Unbalanced number of werewolves",
// 			room: data.WaitingRoom{
// 				PlayerIDs: []types.PlayerID{
// 					"1", "2", "3", "4",
// 				},
// 			},
// 			cfg: data.GameConfig{
// 				NumberWerewolves: 2,
// 			},
// 			expectedErr: fmt.Errorf("Unblanced number of werewolves!"),
// 		},
// 		{
// 			name: "Unbalanced number of werewolves",
// 			room: data.WaitingRoom{
// 				PlayerIDs: []types.PlayerID{
// 					"1", "2", "3", "4", "5",
// 				},
// 			},
// 			cfg: data.GameConfig{
// 				NumberWerewolves: 3,
// 			},
// 			expectedErr: fmt.Errorf("Unblanced number of werewolves!"),
// 		},
// 		{
// 			name: "Ok",
// 			room: data.WaitingRoom{
// 				PlayerIDs: []types.PlayerID{
// 					"1", "2", "3", "4",
// 				},
// 			},
// 			cfg: data.GameConfig{
// 				NumberWerewolves: 1,
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			room: data.WaitingRoom{
// 				PlayerIDs: []types.PlayerID{
// 					"1", "2", "3", "4", "5",
// 				},
// 			},
// 			cfg: data.GameConfig{
// 				NumberWerewolves: 2,
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		gss.Run(test.name, func() {
// 			svc := service.NewGameService(config.Game{
// 				MinCapacity: 4,
// 				MaxCapacity: 10,
// 			}, nil, nil, nil)
// 			err := svc.CheckBeforeRegistration(test.room, test.cfg)

// 			gss.Require().Equal(test.expectedErr, err)
// 		})
// 	}
// }

// func (gss GameServiceSuite) TestRegisterGame() {
// 	tests := []struct {
// 		name        string
// 		config      data.GameConfig
// 		expectedErr error
// 		setup       func(
// 			pdb *mock_storage.MockStore,
// 			mn *mock_game.MockManager,
// 			mod *mock_game.MockModerator,
// 		)
// 	}{
// 		{
// 			name:        "Failure (Store game failed)",
// 			expectedErr: fmt.Errorf("Something went wrong!"),
// 			setup: func(pdb *mock_storage.MockStore, mn *mock_game.MockManager, mod *mock_game.MockModerator) {
// 				pdb.EXPECT().CreateGame(context.Background()).
// 					Return(postgres.Game{}, fmt.Errorf("Store failed"))
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			config: data.GameConfig{
// 				NumberWerewolves: 99,
// 				RoleIDs:          []types.RoleID{vars.HunterRoleID},
// 			},
// 			expectedErr: nil,
// 			setup: func(pdb *mock_storage.MockStore, mn *mock_game.MockManager, mod *mock_game.MockModerator) {
// 				pdb.EXPECT().CreateGame(context.Background()).
// 					Return(postgres.Game{
// 						ID: gss.gameID,
// 					}, nil)
// 				mn.EXPECT().RegisterGame(&types.GameRegistration{
// 					ID: types.GameID(gss.gameID),
// 					GameInitialization: types.GameInitialization{
// 						NumberWerewolves: 99,
// 						RoleIDs:          []types.RoleID{vars.HunterRoleID},
// 						PlayerIDs:        gss.joinedPlayerIDs,
// 					},
// 				}).
// 					Return(mod, nil)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		gss.Run(test.name, func() {
// 			ctrl := gomock.NewController(gss.T())
// 			defer ctrl.Finish()
// 			pdb := mock_storage.NewMockStore(ctrl)
// 			rdb, _ := redismock.NewClusterMock()
// 			moderator := mock_game.NewMockModerator(ctrl)
// 			manager := mock_game.NewMockManager(ctrl)

// 			test.setup(pdb, manager, moderator)

// 			svc := service.NewGameService(config.Game{}, rdb, pdb, manager)
// 			mod, err := svc.RegisterGame(test.config, gss.joinedPlayerIDs)

// 			if test.expectedErr != nil {
// 				gss.Equal(test.expectedErr, err)
// 				gss.Nil(mod)
// 			} else {
// 				gss.Equal(moderator, mod)
// 			}
// 		})
// 	}
// }
