package role

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type WerewolfSuite struct {
	suite.Suite
	playerID types.PlayerID
}

func TestWerewolfSuite(t *testing.T) {
	suite.Run(t, new(WerewolfSuite))
}

func (ws *WerewolfSuite) SetupSuite() {
	ws.playerID = types.PlayerID("1")
}

func (ws WerewolfSuite) TestNewWerewolf() {
	tests := []struct {
		name        string
		expectedErr string
		setup       func(*gamemock.MockGame, *gamemock.MockPoll)
	}{
		{
			name:        "Failure (Poll does not exist)",
			expectedErr: "Poll does not exist ¯\\_(ツ)_/¯",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(nil).Times(1)
			},
		},
		{
			name: "Ok",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(mp).Times(2)
				mp.EXPECT().AddElectors(ws.playerID).Times(1)
				mp.EXPECT().SetWeight(ws.playerID, uint(1)).Times(1)
			},
		},
	}

	for _, test := range tests {
		ws.Run(test.name, func() {
			ctrl := gomock.NewController(ws.T())
			defer ctrl.Finish()
			game := gamemock.NewMockGame(ctrl)
			player := gamemock.NewMockPlayer(ctrl)
			poll := gamemock.NewMockPoll(ctrl)

			game.EXPECT().Player(ws.playerID).Return(player).AnyTimes()
			test.setup(game, poll)

			w, err := NewWerewolf(game, ws.playerID)

			if test.expectedErr != "" {
				ws.Nil(w)
				ws.NotNil(err)
				ws.Equal(test.expectedErr, err.Error())
			} else {
				ws.Nil(err)
				ws.Equal(vars.WerewolfRoleID, w.ID())
				ws.Equal(vars.NightPhaseID, w.(*werewolf).phaseID)
				ws.Equal(vars.WerewolfFactionID, w.FactionID())
				ws.Equal(vars.FirstRound, w.(*werewolf).beginRoundID)
				ws.Equal(player, w.(*werewolf).player)
				ws.Equal(vars.Unlimited, w.ActiveLimit(0))
				ws.Len(w.(*werewolf).abilities, 1)
				ws.Equal(vars.VoteActionID, w.(*werewolf).abilities[0].action.ID())
			}
		})
	}
}

func (ws WerewolfSuite) TestRegisterTurn() {
	ctrl := gomock.NewController(ws.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	player := gamemock.NewMockPlayer(ctrl)
	poll := gamemock.NewMockPoll(ctrl)
	scheduler := gamemock.NewMockScheduler(ctrl)

	player.EXPECT().ID().Return(ws.playerID).Times(1)
	game.EXPECT().Player(ws.playerID).Return(player).Times(1)
	game.EXPECT().Poll(vars.WerewolfFactionID).Return(poll).Times(2)
	poll.EXPECT().AddElectors(ws.playerID).Times(1)
	poll.EXPECT().SetWeight(ws.playerID, uint(1)).Times(1)
	game.EXPECT().Scheduler().Return(scheduler).Times(1)

	w, _ := NewWerewolf(game, ws.playerID)

	scheduler.EXPECT().AddPlayerTurn(&types.NewPlayerTurn{
		PhaseID:      w.(*werewolf).phaseID,
		TurnID:       w.(*werewolf).turnID,
		BeginRoundID: w.(*werewolf).beginRoundID,
		PlayerID:     ws.playerID,
		RoleID:       w.(*werewolf).id,
		ExpiredAfter: vars.Unlimited,
	}).Times(1)

	w.RegisterTurn()
}
