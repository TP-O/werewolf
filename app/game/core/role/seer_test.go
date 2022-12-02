package role

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewSeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)

	playerID := types.PlayerID("1")
	mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)

	seer, _ := NewSeer(mockGame, playerID)

	assert.Equal(t, config.SeerRoleID, seer.ID())
	assert.Equal(t, config.NightPhaseID, seer.PhaseID())
	assert.Equal(t, config.VillagerFactionID, seer.FactionID())
	assert.Equal(t, config.SeerTurnPriority, seer.Priority())
	assert.Equal(t, types.Round(2), seer.BeginRound())
	assert.Equal(t, config.Unlimited, seer.ActiveLimit(config.PredictActionID))
}
