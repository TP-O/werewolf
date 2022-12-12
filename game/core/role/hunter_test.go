package role

import (
	"testing"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type HunterSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	game      *gamemock.MockGame
	player    *gamemock.MockPlayer
	scheduler *gamemock.MockScheduler
	playerID  enum.PlayerID
}

func TestHunterSuite(t *testing.T) {
	suite.Run(t, new(HunterSuite))
}

func (hs *HunterSuite) SetupSuite() {
	hs.playerID = "1"
}

func (hs *HunterSuite) SetupTest() {
	hs.ctrl = gomock.NewController(hs.T())
	hs.game = gamemock.NewMockGame(hs.ctrl)
	hs.player = gamemock.NewMockPlayer(hs.ctrl)
	hs.scheduler = gamemock.NewMockScheduler(hs.ctrl)
	hs.game.EXPECT().Player(hs.playerID).Return(hs.player).AnyTimes()
}

func (hs *HunterSuite) TearDownTest() {
	hs.ctrl.Finish()
}

func (hs *HunterSuite) TestNewHunter() {
	hunter, _ := NewHunter(hs.game, hs.playerID)

	hs.Equal(enum.HunterRoleID, hunter.ID())
	hs.Equal(enum.DayPhaseID, hunter.PhaseID())
	hs.Equal(enum.VillagerFactionID, hunter.FactionID())
	hs.Equal(enum.HunterTurnPriority, hunter.Priority())
	hs.Equal(enum.FirstRound, hunter.BeginRound())
	hs.Equal(enum.ReachedLimit, hunter.ActiveLimit(enum.KillActionID))
}

func (hs *HunterSuite) TestHunterAfterDeath() {
	tests := []struct {
		name  string
		setup func(contract.Role)
	}{
		{
			name: "Ok (Die at unactive phase)",
			setup: func(hunter contract.Role) {
				hs.game.EXPECT().Scheduler().Return(hs.scheduler).Times(1)
				hs.scheduler.EXPECT().PhaseID().Return(enum.NightPhaseID).Times(1)
				hs.game.EXPECT().Scheduler().Return(hs.scheduler).Times(1)
				hs.scheduler.EXPECT().AddTurn(&types.TurnSetting{
					PhaseID:    hunter.PhaseID(),
					RoleID:     hunter.ID(),
					BeginRound: hunter.BeginRound(),
					Position:   enum.SortedPosition,
				}).Times(1)
			},
		},
		{
			name: "Ok (Die at active phase)",
			setup: func(hunter contract.Role) {
				hs.game.EXPECT().Scheduler().Return(hs.scheduler).Times(1)
				hs.scheduler.EXPECT().PhaseID().Return(hunter.PhaseID()).Times(1)
				hs.game.EXPECT().Scheduler().Return(hs.scheduler).Times(1)
				hs.scheduler.EXPECT().AddTurn(&types.TurnSetting{
					PhaseID:    hunter.PhaseID(),
					RoleID:     hunter.ID(),
					BeginRound: hunter.BeginRound(),
					Position:   enum.NextPosition,
				}).Times(1)
			},
		},
	}

	for _, test := range tests {
		hs.Run(test.name, func() {
			hunter, _ := NewHunter(hs.game, hs.playerID)
			test.setup(hunter)
			hunter.AfterDeath()

			hs.Equal(enum.OneMore, hunter.ActiveLimit(enum.KillActionID))
		})
	}
}
