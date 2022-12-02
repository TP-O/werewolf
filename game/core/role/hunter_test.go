package role

import (
	"testing"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewHunter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)

	playerID := enum.PlayerID("1")
	mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)

	hunter, _ := NewHunter(mockGame, playerID)

	assert.Equal(t, enum.HunterRoleID, hunter.ID())
	assert.Equal(t, enum.DayPhaseID, hunter.PhaseID())
	assert.Equal(t, enum.VillagerFactionID, hunter.FactionID())
	assert.Equal(t, enum.HunterTurnPriority, hunter.Priority())
	assert.Equal(t, enum.FirstRound, hunter.BeginRound())
	assert.Equal(t, enum.ReachedLimit, hunter.ActiveLimit(enum.KillActionID))
}

func TestHunterAfterDeath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)
	mockScheduler := gamemock.NewMockScheduler(ctrl)

	playerID := enum.PlayerID("1")

	tests := []struct {
		name  string
		setup func(contract.Role)
	}{
		{
			name: "Ok (Die at unactive phase)",
			setup: func(hunter contract.Role) {
				mockGame.EXPECT().Scheduler().Return(mockScheduler).Times(1)
				mockScheduler.EXPECT().PhaseID().Return(enum.NightPhaseID).Times(1)
				mockGame.EXPECT().Scheduler().Return(mockScheduler).Times(1)
				mockScheduler.EXPECT().AddTurn(&types.TurnSetting{
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
				mockGame.EXPECT().Scheduler().Return(mockScheduler).Times(1)
				mockScheduler.EXPECT().PhaseID().Return(hunter.PhaseID()).Times(1)
				mockGame.EXPECT().Scheduler().Return(mockScheduler).Times(1)
				mockScheduler.EXPECT().AddTurn(&types.TurnSetting{
					PhaseID:    hunter.PhaseID(),
					RoleID:     hunter.ID(),
					BeginRound: hunter.BeginRound(),
					Position:   enum.NextPosition,
				}).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)
			hunter, _ := NewHunter(mockGame, playerID)
			test.setup(hunter)
			hunter.AfterDeath()

			assert.Equal(t, enum.OneMore, hunter.ActiveLimit(enum.KillActionID))
		})
	}
}
