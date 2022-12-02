package role

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTwoSister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)

	playerID := types.PlayerID("1")
	mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)

	twoSister, _ := NewTwoSister(mockGame, playerID)

	assert.Equal(t, config.TwoSistersRoleID, twoSister.ID())
	assert.Equal(t, config.NightPhaseID, twoSister.PhaseID())
	assert.Equal(t, config.VillagerFactionID, twoSister.FactionID())
	assert.Equal(t, config.TwoSistersTurnPriority, twoSister.Priority())
	assert.Equal(t, config.FirstRound, twoSister.BeginRound())
	assert.Equal(t, config.OneMore, twoSister.ActiveLimit(config.RecognizeActionID))
}
