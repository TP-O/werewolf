package role

import (
	"testing"
	"uwwolf/game/enum"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTwoSister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)

	playerID := enum.PlayerID("1")
	mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)

	twoSister, _ := NewTwoSister(mockGame, playerID)

	assert.Equal(t, enum.TwoSistersRoleID, twoSister.ID())
	assert.Equal(t, enum.NightPhaseID, twoSister.PhaseID())
	assert.Equal(t, enum.VillagerFactionID, twoSister.FactionID())
	assert.Equal(t, enum.TwoSistersTurnPriority, twoSister.Priority())
	assert.Equal(t, enum.FirstRound, twoSister.BeginRound())
	assert.Equal(t, enum.OneMore, twoSister.ActiveLimit(enum.RecognizeActionID))
}
