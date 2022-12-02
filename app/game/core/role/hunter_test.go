package role

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewHunter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)

	playerID := types.PlayerID("1")
	mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)

	hunter, _ := NewHunter(mockGame, playerID)

	assert.Equal(t, config.HunterRoleID, hunter.ID())
	assert.Equal(t, config.DayPhaseID, hunter.PhaseID())
	assert.Equal(t, config.VillagerFactionID, hunter.FactionID())
	assert.Equal(t, config.HunterTurnPriority, hunter.Priority())
	assert.Equal(t, config.FirstRound, hunter.BeginRound())
	assert.Equal(t, config.ReachedLimit, hunter.ActiveLimit(config.KillActionID))
}

func TestHunterAfterDeath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)
	mockScheduler := gamemock.NewMockScheduler(ctrl)

	playerID := types.PlayerID("1")

	tests := []struct {
		name  string
		setup func(contract.Role)
	}{
		{
			name: "Die at unactive phase",
			setup: func(hunter contract.Role) {
				mockGame.EXPECT().Scheduler().Return(mockScheduler).Times(1)
				mockScheduler.EXPECT().PhaseID().Return(config.NightPhaseID).Times(1)
				mockGame.EXPECT().Scheduler().Return(mockScheduler).Times(1)
				mockScheduler.EXPECT().AddTurn(&types.TurnSetting{
					PhaseID:    hunter.PhaseID(),
					RoleID:     hunter.ID(),
					BeginRound: hunter.BeginRound(),
					Position:   config.SortedPosition,
				}).Times(1)
			},
		},
		{
			name: "Die at active phase",
			setup: func(hunter contract.Role) {
				mockGame.EXPECT().Scheduler().Return(mockScheduler).Times(1)
				mockScheduler.EXPECT().PhaseID().Return(hunter.PhaseID()).Times(1)
				mockGame.EXPECT().Scheduler().Return(mockScheduler).Times(1)
				mockScheduler.EXPECT().AddTurn(&types.TurnSetting{
					PhaseID:    hunter.PhaseID(),
					RoleID:     hunter.ID(),
					BeginRound: hunter.BeginRound(),
					Position:   config.NextPosition,
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

			assert.Equal(t, config.OneMore, hunter.ActiveLimit(config.KillActionID))
		})
	}
}
