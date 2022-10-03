package action_test

// import (
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"

// 	"uwwolf/app/game/action"
// 	"uwwolf/mock/game"
// 	"uwwolf/types"
// )

// func TestShootingId(t *testing.T) {
// 	p := action.NewShooting(nil)

// 	assert.Equal(t, types.ShootingAction, p.Id())
// }

// func TestShootingState(t *testing.T) {
// 	p := action.NewShooting(nil)

// 	assert.Nil(t, p.State())
// }

// func TestShootingPerform(t *testing.T) {
// 	//========================MOCK================================
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockGame := game.NewMockGame(ctrl)
// 	mockPlayer := game.NewMockPlayer(ctrl)

// 	//=============================================================
// 	hunterId := types.PlayerId("1")
// 	targetId := types.PlayerId("2")
// 	mockGame.
// 		EXPECT().
// 		KillPlayer(gomock.Any()).
// 		Return(mockPlayer).
// 		Times(2)

// 	mockPlayer.
// 		EXPECT().
// 		Id().
// 		Return(targetId).
// 		Times(2)

// 	//=============================================================
// 	// ActorId and target is the same
// 	p := action.NewShooting(mockGame)

// 	res := p.Perform(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   hunterId,
// 		TargetIds: []types.PlayerId{hunterId},
// 	})

// 	assert.False(t, res.Ok)

// 	//=============================================================
// 	// Already shot
// 	p = action.NewShooting(mockGame)

// 	p.Perform(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   hunterId,
// 		TargetIds: []types.PlayerId{targetId},
// 	})

// 	res = p.Perform(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   hunterId,
// 		TargetIds: []types.PlayerId{targetId},
// 	})

// 	assert.False(t, res.Ok)

// 	//=============================================================
// 	// Skip shooting
// 	p = action.NewShooting(mockGame)

// 	res = p.Perform(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   hunterId,
// 		IsSkipped: true,
// 	})

// 	assert.True(t, res.Ok)

// 	//=============================================================
// 	// Successfully shot
// 	res = p.Perform(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   hunterId,
// 		TargetIds: []types.PlayerId{targetId},
// 	})

// 	assert.True(t, res.Ok)
// 	assert.Equal(t, targetId, res.Data)
// }
