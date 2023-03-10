package role

import (
	"fmt"
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
		expectedErr error
		setup       func(*gamemock.MockGame, *gamemock.MockPoll)
	}{
		{
			name:        "Failure (Poll does not exist)",
			expectedErr: fmt.Errorf("Poll does not exist ¯\\_(ツ)_/¯"),
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(nil)
			},
		},
		{
			name: "Ok",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(mp).Times(2)
				mp.EXPECT().AddElectors(ws.playerID)
				mp.EXPECT().SetWeight(ws.playerID, uint(1))
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

			if test.expectedErr != nil {
				ws.Nil(w)
				ws.NotNil(err)
				ws.Equal(test.expectedErr, err)
			} else {
				ws.Nil(err)
				ws.Equal(vars.WerewolfRoleID, w.ID())
				ws.Equal(vars.NightPhaseID, w.(*werewolf).phaseID)
				ws.Equal(vars.WerewolfFactionID, w.FactionID())
				ws.Equal(vars.FirstRound, w.(*werewolf).beginRoundID)
				ws.Equal(player, w.(*werewolf).player)
				ws.Equal(vars.UnlimitedTimes, w.ActiveTimes(0))
				ws.Len(w.(*werewolf).abilities, 1)
				ws.Equal(vars.VoteActionID, w.(*werewolf).abilities[0].action.ID())
			}
		})
	}
}

func (ws WerewolfSuite) TestOnRevoke() {
	ctrl := gomock.NewController(ws.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	player := gamemock.NewMockPlayer(ctrl)
	scheduler := gamemock.NewMockScheduler(ctrl)
	poll := gamemock.NewMockPoll(ctrl)

	// Mock for New fuction
	game.EXPECT().Scheduler().Return(scheduler)
	game.EXPECT().Player(ws.playerID).Return(player)
	game.EXPECT().Poll(vars.WerewolfFactionID).Return(poll).Times(2)
	poll.EXPECT().AddElectors(ws.playerID)
	poll.EXPECT().SetWeight(ws.playerID, uint(1))

	game.EXPECT().Poll(vars.VillagerFactionID).Return(poll).Times(2)
	game.EXPECT().Poll(vars.WerewolfFactionID).Return(poll)
	poll.EXPECT().RemoveElector(ws.playerID).Times(2)
	poll.EXPECT().RemoveCandidate(ws.playerID)
	player.EXPECT().ID().Return(ws.playerID).Times(4)

	w, _ := NewWerewolf(game, ws.playerID)

	scheduler.EXPECT().RemoveSlot(&types.RemovedTurnSlot{
		PhaseID:  w.(*werewolf).phaseID,
		PlayerID: ws.playerID,
		RoleID:   w.ID(),
	})

	w.OnRevoke()
}

func (ws WerewolfSuite) TestOnAssign() {
	ctrl := gomock.NewController(ws.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	player := gamemock.NewMockPlayer(ctrl)
	scheduler := gamemock.NewMockScheduler(ctrl)
	poll := gamemock.NewMockPoll(ctrl)

	// Mock for New fuction
	game.EXPECT().Scheduler().Return(scheduler)
	game.EXPECT().Player(ws.playerID).Return(player)
	game.EXPECT().Poll(vars.WerewolfFactionID).Return(poll).Times(2)
	poll.EXPECT().AddElectors(ws.playerID)
	poll.EXPECT().SetWeight(ws.playerID, uint(1))

	game.EXPECT().Poll(vars.VillagerFactionID).Return(poll)
	poll.EXPECT().AddCandidates(ws.playerID)
	player.EXPECT().ID().Return(ws.playerID).Times(2)

	w, _ := NewWerewolf(game, ws.playerID)

	scheduler.EXPECT().AddSlot(&types.NewTurnSlot{
		PhaseID:      w.(*werewolf).phaseID,
		TurnID:       w.(*werewolf).turnID,
		BeginRoundID: w.(*werewolf).beginRoundID,
		PlayerID:     ws.playerID,
		RoleID:       w.ID(),
	})

	w.OnAssign()
}
