package role

import (
	"testing"
	"uwwolf/game/enum"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type WerewolfSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	game     *gamemock.MockGame
	player   *gamemock.MockPlayer
	poll     *gamemock.MockPoll
	playerID enum.PlayerID
}

func TestWerewolfSuite(t *testing.T) {
	suite.Run(t, new(WerewolfSuite))
}

func (ws *WerewolfSuite) SetupSuite() {
	ws.playerID = "1"
}

func (ws *WerewolfSuite) SetupTest() {
	ws.ctrl = gomock.NewController(ws.T())
	ws.game = gamemock.NewMockGame(ws.ctrl)
	ws.player = gamemock.NewMockPlayer(ws.ctrl)
	ws.poll = gamemock.NewMockPoll(ws.ctrl)
	ws.game.EXPECT().Player(ws.playerID).Return(ws.player).AnyTimes()
}

func (ws *WerewolfSuite) TearDownTest() {
	ws.ctrl.Finish()
}

func (ws *WerewolfSuite) TestNewWerewolf() {
	tests := []struct {
		name        string
		returnNil   bool
		expectedRrr string
		setup       func()
	}{
		{
			name:        "Failure (Poll does not exist)",
			returnNil:   true,
			expectedRrr: "Poll does not exist ¯\\_(ツ)_/¯",
			setup: func() {
				ws.game.EXPECT().Poll(enum.WerewolfFactionID).Return(nil).Times(1)
			},
		},
		{
			name:        "Failure (Unable to join to the poll)",
			returnNil:   true,
			expectedRrr: "Unable to join to the poll ಠ_ಠ",
			setup: func() {
				ws.game.EXPECT().Poll(enum.WerewolfFactionID).Return(ws.poll).Times(1)
				ws.poll.EXPECT().AddElectors(ws.playerID).Return(false).Times(1)
			},
		},
		{
			name:      "Ok",
			returnNil: false,
			setup: func() {
				ws.game.EXPECT().Poll(enum.WerewolfFactionID).Return(ws.poll).Times(1)
				ws.poll.EXPECT().AddElectors(ws.playerID).Return(true).Times(1)
				ws.poll.EXPECT().SetWeight(ws.playerID, uint(1)).Times(1)
			},
		},
	}

	for _, test := range tests {
		ws.Run(test.name, func() {
			test.setup()
			werewolf, err := NewWerewolf(ws.game, ws.playerID)

			if test.returnNil {
				ws.Nil(werewolf)
				ws.NotNil(err)
				ws.Equal(test.expectedRrr, err.Error())
			} else {
				ws.Nil(err)
				ws.Equal(enum.WerewolfRoleID, werewolf.ID())
				ws.Equal(enum.NightPhaseID, werewolf.PhaseID())
				ws.Equal(enum.WerewolfFactionID, werewolf.FactionID())
				ws.Equal(enum.WerewolfTurnPriority, werewolf.Priority())
				ws.Equal(enum.FirstRound, werewolf.BeginRound())
				ws.Equal(enum.Unlimited, werewolf.ActiveLimit(enum.VoteActionID))
			}
		})
	}
}
