package game

// import (
// 	"testing"
// 	"time"
// 	"uwwolf/config"
// 	"uwwolf/game/types"
// 	"uwwolf/game/vars"
// 	mock_game "uwwolf/mock/game"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type ModeratorSuite struct {
// 	suite.Suite
// 	playerID types.PlayerID
// }

// func (ms *ModeratorSuite) SetupSuite() {
// 	ms.playerID = types.PlayerID("1")
// }

// func TestModeratorSuite(t *testing.T) {
// 	suite.Run(t, new(ModeratorSuite))
// }

// func (ms ModeratorSuite) TestNewModerator() {
// 	ctrl := gomock.NewController(ms.T())
// 	defer ctrl.Finish()

// 	reg := &types.GameRegistration{
// 		ID: types.GameID(1),

// 		TurnDuration:       5 * time.Second,
// 		DiscussionDuration: 10 * time.Second,
// 		GameInitialization: types.GameInitialization{},
// 	}
// 	m := NewModerator(config.Game{}, reg).(*moderator)

// 	ms.Equal(reg.ID, m.gameID)
// 	ms.Equal(reg.TurnDuration, m.turnDuration)
// 	ms.Equal(reg.DiscussionDuration, m.discussionDuration)
// 	ms.NotNil(m.game)
// 	ms.NotNil(m.scheduler)
// 	ms.NotNil(m.finishSignal)
// 	ms.NotNil(m.nextTurnSignal)
// 	ms.NotNil(m.mutex)
// }

// func (ms ModeratorSuite) TestCheckWinConditions() {
// 	tests := []struct {
// 		name                     string
// 		expectedWinningFactionID types.FactionID
// 		setup                    func(*moderator, *mock_game.MockGame)
// 	}{
// 		{
// 			name:                     "Villager wins",
// 			expectedWinningFactionID: vars.VillagerFactionID,
// 			setup: func(m *moderator, mg *mock_game.MockGame) {
// 				mg.EXPECT().AlivePlayerIDsWithFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{})
// 				mg.EXPECT().StatusID().Return(vars.Waiting)
// 				mg.EXPECT().Finish()
// 				go func() {
// 					<-m.finishSignal
// 				}()
// 			},
// 		},
// 		{
// 			name:                     "Werewolf wins",
// 			expectedWinningFactionID: vars.WerewolfFactionID,
// 			setup: func(m *moderator, mg *mock_game.MockGame) {
// 				mg.EXPECT().AlivePlayerIDsWithFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1", "2"}).Times(2)
// 				mg.EXPECT().AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"3", "4"})
// 				mg.EXPECT().StatusID().Return(vars.Waiting)
// 				mg.EXPECT().Finish()
// 				go func() {
// 					<-m.finishSignal
// 				}()
// 			},
// 		},
// 		{
// 			name:                     "Neither faction wins",
// 			expectedWinningFactionID: types.FactionID(0),
// 			setup: func(m *moderator, mg *mock_game.MockGame) {
// 				mg.EXPECT().AlivePlayerIDsWithFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1", "2"}).Times(2)
// 				mg.EXPECT().AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"3", "4", "5"})
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ms.Run(test.name, func() {
// 			ctrl := gomock.NewController(ms.T())
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)

// 			m := NewModerator(config.Game{}, &types.GameRegistration{}).(*moderator)
// 			m.game = game
// 			test.setup(m, game)

// 			m.checkWinConditions()

// 			ms.Equal(test.expectedWinningFactionID, m.winningFaction)
// 		})
// 	}
// }

// func (ms ModeratorSuite) TestHandlePoll() {
// 	tests := []struct {
// 		name      string
// 		factionID types.FactionID
// 		setup     func(*mock_game.MockGame, *mock_game.MockPoll)
// 	}{
// 		{
// 			name:      "Unsupported poll",
// 			factionID: types.FactionID(99),
// 			setup: func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {
// 				mg.EXPECT().Poll(types.FactionID(99)).Return(nil)
// 			},
// 		},
// 		{
// 			name:      "Villager poll",
// 			factionID: vars.VillagerFactionID,
// 			setup: func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {
// 				mg.EXPECT().Poll(vars.VillagerFactionID).Return(mp)
// 				mp.EXPECT().Close().Return(true)
// 				mp.EXPECT().Record(vars.ZeroRound).Return(&types.PollRecord{
// 					WinnerID: ms.playerID,
// 				})
// 				mg.EXPECT().KillPlayer(ms.playerID, false)
// 			},
// 		},
// 		{
// 			name:      "Werewolf poll",
// 			factionID: vars.WerewolfFactionID,
// 			setup: func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {
// 				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(mp)
// 				mp.EXPECT().Close().Return(true)
// 				mp.EXPECT().Record(vars.ZeroRound).Return(&types.PollRecord{
// 					WinnerID: ms.playerID,
// 				})
// 				mg.EXPECT().KillPlayer(ms.playerID, false)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ms.Run(test.name, func() {
// 			ctrl := gomock.NewController(ms.T())
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)
// 			poll := mock_game.NewMockPoll(ctrl)

// 			m := NewModerator(config.Game{}, &types.GameRegistration{}).(*moderator)
// 			m.game = game
// 			test.setup(game, poll)

// 			m.handlePoll(test.factionID)
// 		})
// 	}
// }

// func (ms ModeratorSuite) TestRunScheduler() {
// 	tests := []struct {
// 		name  string
// 		setup func(*moderator, *mock_game.MockGame, *mock_game.MockScheduler, *mock_game.MockPoll)
// 	}{
// 		{
// 			name: "Game isn't starting",
// 			setup: func(m *moderator, mg *mock_game.MockGame, ms *mock_game.MockScheduler, mp *mock_game.MockPoll) {
// 				mg.EXPECT().StatusID().Return(vars.Waiting)
// 			},
// 		},
// 		{
// 			name: "Villager turn",
// 			setup: func(m *moderator, mg *mock_game.MockGame, ms *mock_game.MockScheduler, mp *mock_game.MockPoll) {
// 				m.turnDuration = 99 * time.Hour
// 				m.discussionDuration = 0 * time.Second

// 				mg.EXPECT().StatusID().Return(vars.Starting)
// 				ms.EXPECT().NextTurn()
// 				ms.EXPECT().PhaseID().Return(vars.DayPhaseID)
// 				ms.EXPECT().TurnID().Return(vars.MidTurn)
// 				mg.EXPECT().Poll(vars.VillagerFactionID).Return(mp)
// 				mp.EXPECT().Open()
// 				mg.EXPECT().Poll(vars.VillagerFactionID).Return(nil)
// 				mg.EXPECT().AlivePlayerIDsWithFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1"}).Times(2)
// 				mg.EXPECT().AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1", "2"})
// 				mg.EXPECT().StatusID().Return(vars.Finished)
// 			},
// 		},
// 		{
// 			name: "Werewolf turn",
// 			setup: func(m *moderator, mg *mock_game.MockGame, ms *mock_game.MockScheduler, mp *mock_game.MockPoll) {
// 				m.turnDuration = 0 * time.Second
// 				m.discussionDuration = 99 * time.Hour

// 				mg.EXPECT().StatusID().Return(vars.Starting)
// 				ms.EXPECT().NextTurn()
// 				ms.EXPECT().PhaseID().Return(vars.NightPhaseID).Times(2)
// 				ms.EXPECT().TurnID().Return(vars.MidTurn)
// 				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(mp)
// 				mp.EXPECT().Open()
// 				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(nil)
// 				mg.EXPECT().AlivePlayerIDsWithFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1"}).Times(2)
// 				mg.EXPECT().AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1", "2"})
// 				mg.EXPECT().StatusID().Return(vars.Finished)
// 			},
// 		},
// 		{
// 			name: "Normal turn",
// 			setup: func(m *moderator, mg *mock_game.MockGame, ms *mock_game.MockScheduler, mp *mock_game.MockPoll) {
// 				m.turnDuration = 0 * time.Second
// 				m.discussionDuration = 99 * time.Hour

// 				mg.EXPECT().StatusID().Return(vars.Starting)
// 				ms.EXPECT().NextTurn()
// 				ms.EXPECT().PhaseID().Return(vars.NightPhaseID).Times(2)
// 				ms.EXPECT().TurnID().Return(vars.PreTurn)
// 				mg.EXPECT().AlivePlayerIDsWithFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1"}).Times(2)
// 				mg.EXPECT().AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1", "2"})
// 				mg.EXPECT().StatusID().Return(vars.Finished)
// 			},
// 		},
// 		{
// 			name: "Next turn by nextTurnSignal",
// 			setup: func(m *moderator, mg *mock_game.MockGame, ms *mock_game.MockScheduler, mp *mock_game.MockPoll) {
// 				m.turnDuration = 99 * time.Hour
// 				m.discussionDuration = 99 * time.Hour

// 				mg.EXPECT().StatusID().Return(vars.Starting)
// 				ms.EXPECT().NextTurn()
// 				ms.EXPECT().PhaseID().Return(vars.DayPhaseID).Times(2)
// 				ms.EXPECT().TurnID().Return(vars.PreTurn)
// 				mg.EXPECT().AlivePlayerIDsWithFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1"}).Times(2)
// 				mg.EXPECT().AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1", "2"})
// 				mg.EXPECT().StatusID().Return(vars.Finished)

// 				go func() {
// 					m.nextTurnSignal <- true
// 				}()
// 			},
// 		},
// 		{
// 			name: "Next turn by nextTurnSignal",
// 			setup: func(m *moderator, mg *mock_game.MockGame, ms *mock_game.MockScheduler, mp *mock_game.MockPoll) {
// 				m.turnDuration = 99 * time.Hour
// 				m.discussionDuration = 99 * time.Hour

// 				mg.EXPECT().StatusID().Return(vars.Starting)
// 				ms.EXPECT().NextTurn()
// 				ms.EXPECT().PhaseID().Return(vars.DayPhaseID).Times(2)
// 				ms.EXPECT().TurnID().Return(vars.PreTurn)
// 				mg.EXPECT().AlivePlayerIDsWithFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1"}).Times(2)
// 				mg.EXPECT().AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID).
// 					Return([]types.PlayerID{"1", "2"})
// 				mg.EXPECT().StatusID().Return(vars.Finished)

// 				go func() {
// 					m.nextTurnSignal <- true
// 				}()
// 			},
// 		},
// 		{
// 			name: "Next turn by finishSignal",
// 			setup: func(m *moderator, mg *mock_game.MockGame, ms *mock_game.MockScheduler, mp *mock_game.MockPoll) {
// 				m.turnDuration = 99 * time.Hour
// 				m.discussionDuration = 99 * time.Hour

// 				mg.EXPECT().StatusID().Return(vars.Starting)
// 				ms.EXPECT().NextTurn()
// 				ms.EXPECT().PhaseID().Return(vars.DayPhaseID).Times(2)
// 				ms.EXPECT().TurnID().Return(vars.PreTurn)
// 				mg.EXPECT().StatusID().Return(vars.Finished).Times(2)

// 				go func() {
// 					m.finishSignal <- true
// 				}()
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ms.Run(test.name, func() {
// 			ctrl := gomock.NewController(ms.T())
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)
// 			scheduler := mock_game.NewMockScheduler(ctrl)
// 			poll := mock_game.NewMockPoll(ctrl)

// 			m := NewModerator(config.Game{}, &types.GameRegistration{}).(*moderator)
// 			m.game = game
// 			m.scheduler = scheduler
// 			test.setup(m, game, scheduler, poll)

// 			m.runScheduler()
// 		})
// 	}
// }

// func (ms ModeratorSuite) TestWaitForPreparation() {
// 	tests := []struct {
// 		name   string
// 		config config.Game
// 		setup  func(mod *moderator, game *mock_game.MockGame)
// 	}{
// 		{
// 			name: "Wait until timeout",
// 			config: config.Game{
// 				PreparationDuration: 0,
// 			},
// 			setup: func(mod *moderator, game *mock_game.MockGame) {},
// 		},
// 		{
// 			name: "Wait until finishSignal emitted",
// 			config: config.Game{
// 				PreparationDuration: 99 * time.Second,
// 			},
// 			setup: func(mod *moderator, game *mock_game.MockGame) {
// 				game.EXPECT().StatusID().Return(vars.Finished)

// 				go func() {
// 					mod.finishSignal <- true
// 				}()
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ms.Run(test.name, func() {
// 			ctrl := gomock.NewController(ms.T())
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)

// 			m := NewModerator(test.config, &types.GameRegistration{}).(*moderator)
// 			m.game = game
// 			test.setup(m, game)

// 			m.waitForPreparation()
// 		})
// 	}
// }

// func (ms ModeratorSuite) TestStartGame() {
// 	tests := []struct {
// 		name         string
// 		expectedTime int64
// 		setup        func(game *mock_game.MockGame)
// 	}{
// 		{
// 			name:         "False (Game status isn't idle)",
// 			expectedTime: -1,
// 			setup: func(game *mock_game.MockGame) {
// 				game.EXPECT().StatusID().Return(vars.Starting)
// 			},
// 		},
// 		{
// 			name:         "False (Game preparation is failed)",
// 			expectedTime: -1,
// 			setup: func(game *mock_game.MockGame) {
// 				game.EXPECT().StatusID().Return(vars.Idle)
// 				game.EXPECT().Start().MaxTimes(1)
// 				game.EXPECT().Prepare().Return(int64(-1))
// 				game.EXPECT().StatusID().Return(vars.Finished).MaxTimes(1)
// 			},
// 		},
// 		{
// 			name:         "Ok",
// 			expectedTime: 999,
// 			setup: func(game *mock_game.MockGame) {
// 				game.EXPECT().StatusID().Return(vars.Idle)
// 				game.EXPECT().Prepare().Return(int64(999))
// 				game.EXPECT().Start().MaxTimes(1)
// 				game.EXPECT().StatusID().Return(vars.Finished).MaxTimes(1)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ms.Run(test.name, func() {
// 			ctrl := gomock.NewController(ms.T())
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)

// 			m := NewModerator(config.Game{}, &types.GameRegistration{
// 				ID: types.GameID(1),
// 			}).(*moderator)
// 			m.game = game

// 			test.setup(game)

// 			t := m.StartGame()

// 			ms.Equal(test.expectedTime, t)
// 		})
// 	}
// }

// func (ms ModeratorSuite) TestFinishGame() {
// 	tests := []struct {
// 		name           string
// 		expectedStatus bool
// 		setup          func(*moderator, *mock_game.MockGame)
// 	}{
// 		{
// 			name:           "False (Game status is already finished)",
// 			expectedStatus: false,
// 			setup: func(m *moderator, mg *mock_game.MockGame) {
// 				mg.EXPECT().StatusID().Return(vars.Finished)
// 			},
// 		},
// 		{
// 			name:           "Ok",
// 			expectedStatus: true,
// 			setup: func(m *moderator, mg *mock_game.MockGame) {
// 				mg.EXPECT().StatusID().Return(vars.Starting)
// 				mg.EXPECT().Finish()

// 				go func() {
// 					<-m.finishSignal
// 				}()
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ms.Run(test.name, func() {
// 			ctrl := gomock.NewController(ms.T())
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)

// 			m := NewModerator(config.Game{}, &types.GameRegistration{}).(*moderator)
// 			m.game = game
// 			test.setup(m, game)

// 			ok := m.FinishGame()

// 			ms.Equal(test.expectedStatus, ok)
// 			if test.expectedStatus == true {
// 				_, ok1 := <-m.finishSignal
// 				_, ok2 := <-m.nextTurnSignal
// 				ms.False(ok1)
// 				ms.False(ok2)
// 			}
// 		})
// 	}
// }

// func (ms ModeratorSuite) TestRequestPlay() {
// 	tests := []struct {
// 		name        string
// 		expectedRes *types.ActionResponse
// 		setup       func(*moderator, *mock_game.MockGame)
// 	}{
// 		{
// 			name: "Failure (Play is locked)",
// 			expectedRes: &types.ActionResponse{
// 				Ok:      false,
// 				Message: "Turn is over!",
// 			},
// 			setup: func(m *moderator, mg *mock_game.MockGame) {
// 				m.mutex.Lock()
// 			},
// 		},
// 		{
// 			name: "Failure (Play is already played)",
// 			expectedRes: &types.ActionResponse{
// 				Ok:      false,
// 				Message: "You played this turn!",
// 			},
// 			setup: func(m *moderator, mg *mock_game.MockGame) {
// 				m.playedPlayerID = append(m.playedPlayerID, ms.playerID)
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			expectedRes: &types.ActionResponse{
// 				Ok: true,
// 			},
// 			setup: func(m *moderator, mg *mock_game.MockGame) {
// 				mg.EXPECT().Play(ms.playerID, &types.ActivateAbilityRequest{}).
// 					Return(&types.ActionResponse{
// 						Ok: true,
// 					})
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ms.Run(test.name, func() {
// 			ctrl := gomock.NewController(ms.T())
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)

// 			m := NewModerator(config.Game{}, &types.GameRegistration{}).(*moderator)
// 			m.game = game
// 			test.setup(m, game)

// 			res := m.RequestPlay(ms.playerID, &types.ActivateAbilityRequest{})

// 			ms.Equal(test.expectedRes, res)
// 			if test.expectedRes.Ok {
// 				ms.Contains(m.playedPlayerID, ms.playerID)
// 			}
// 		})
// 	}
// }
