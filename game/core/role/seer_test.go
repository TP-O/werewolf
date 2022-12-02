package role

import (
	"testing"
	"uwwolf/game/enum"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewSeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)

	playerID := enum.PlayerID("1")
	mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)

	seer, _ := NewSeer(mockGame, playerID)

	assert.Equal(t, enum.SeerRoleID, seer.ID())
	assert.Equal(t, enum.NightPhaseID, seer.PhaseID())
	assert.Equal(t, enum.VillagerFactionID, seer.FactionID())
	assert.Equal(t, enum.SeerTurnPriority, seer.Priority())
	assert.Equal(t, enum.Round(2), seer.BeginRound())
	assert.Equal(t, enum.Unlimited, seer.ActiveLimit(enum.PredictActionID))
}
