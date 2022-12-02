package action_test

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/core/action"
	"uwwolf/app/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPoll := gamemock.NewMockPoll(ctrl)

	factionID := config.WerewolfFactionID
	playerID := types.PlayerID("1")
	weight := uint(5)
	tests := []struct {
		name        string
		setting     *types.VoteActionSetting
		expectedErr string
		setup       func()
	}{
		{
			name: "Poll does not exist",
			setting: &types.VoteActionSetting{
				FactionID: factionID,
				PlayerID:  playerID,
				Weight:    weight,
			},
			expectedErr: "Poll does not exist ¯\\_(ツ)_/¯",
			setup: func() {
				mockGame.EXPECT().Poll(factionID).Return(nil).Times(1)
			},
		},
		{
			name: "Cannot add player as elector to poll",
			setting: &types.VoteActionSetting{
				FactionID: factionID,
				PlayerID:  playerID,
				Weight:    weight,
			},
			expectedErr: "Unable to join to the poll ಠ_ಠ",
			setup: func() {
				mockGame.EXPECT().Poll(factionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(false).Times(1)
			},
		},
		{
			name: "Ok",
			setting: &types.VoteActionSetting{
				FactionID: factionID,
				PlayerID:  playerID,
				Weight:    weight,
			},
			expectedErr: "",
			setup: func() {
				mockGame.EXPECT().Poll(factionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(true).Times(1)
				mockPoll.EXPECT().SetWeight(playerID, weight).Return(true).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			vote, err := action.NewVote(mockGame, test.setting)

			if test.expectedErr != "" {
				assert.Nil(t, vote)
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.Equal(t, config.VoteActionID, vote.ID())
				assert.NotNil(t, vote.State())
			}
		})
	}
}

func TestValidateVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPoll := gamemock.NewMockPoll(ctrl)

	factionID := config.WerewolfFactionID
	actorID := types.PlayerID("1")
	targetID := types.PlayerID("2")
	weight := uint(1)
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr string
		setup       func()
	}{
		{
			name:        "Empty action request",
			req:         nil,
			expectedErr: "Action request can not be empty (⊙＿⊙')",
			setup:       func() {},
		},
		{
			name: "Not allowed to vote",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
			},
			expectedErr: "Not allowed to vote ¯\\_(ツ)_/¯",
			setup: func() {
				mockPoll.EXPECT().CanVote(actorID).Times(1).Return(false).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock for vote action initialization
			mockGame.EXPECT().Poll(factionID).Return(mockPoll).Times(1)
			mockPoll.EXPECT().AddElectors(actorID).Return(true).Times(1)
			mockPoll.EXPECT().SetWeight(actorID, weight).Return(true).Times(1)

			vote, _ := action.NewVote(mockGame, &types.VoteActionSetting{
				FactionID: factionID,
				PlayerID:  actorID,
				Weight:    weight,
			})
			test.setup()
			err := vote.Validate(test.req)

			assert.Equal(t, test.expectedErr, err.Error())
		})
	}
}

func TestSkipVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPoll := gamemock.NewMockPoll(ctrl)

	factionID := config.WerewolfFactionID
	actorID := types.PlayerID("1")
	weight := uint(1)

	// Mock for vote action initialization
	mockGame.EXPECT().Poll(factionID).Return(mockPoll).Times(1)
	mockPoll.EXPECT().AddElectors(actorID).Return(true).Times(1)
	mockPoll.EXPECT().SetWeight(actorID, weight).Return(true).Times(1)
	mockPoll.EXPECT().Vote(actorID, types.PlayerID("")).Times(1)

	expectedRes := &types.ActionResponse{
		Ok:        true,
		IsSkipped: true,
		Data:      nil,
		Message:   "",
	}
	vote, _ := action.NewVote(mockGame, &types.VoteActionSetting{
		FactionID: factionID,
		PlayerID:  actorID,
		Weight:    weight,
	})
	res := vote.Skip(&types.ActionRequest{
		ActorID:   actorID,
		IsSkipped: true,
	})

	assert.Equal(t, expectedRes, res)
}

func TestPerformVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPoll := gamemock.NewMockPoll(ctrl)

	factionID := config.WerewolfFactionID
	actorID := types.PlayerID("1")
	targetID := types.PlayerID("2")
	weight := uint(1)
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
		setup       func()
	}{

		{
			name: "Vote failed",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        false,
				IsSkipped: false,
				Data:      nil,
				Message:   "Unable to vote (╯°□°)╯︵ ┻━┻",
			},
			setup: func() {
				mockPoll.EXPECT().Vote(actorID, targetID).Times(1).Return(false).Times(1)
			},
		},
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      targetID,
				Message:   "",
			},
			setup: func() {
				mockPoll.EXPECT().Vote(actorID, targetID).Times(1).Return(true).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock for vote action initialization
			mockGame.EXPECT().Poll(factionID).Return(mockPoll).Times(1)
			mockPoll.EXPECT().AddElectors(actorID).Return(true).Times(1)
			mockPoll.EXPECT().SetWeight(actorID, weight).Return(true).Times(1)

			vote, _ := action.NewVote(mockGame, &types.VoteActionSetting{
				FactionID: factionID,
				PlayerID:  actorID,
				Weight:    weight,
			})
			test.setup()
			res := vote.Perform(test.req)

			assert.Equal(t, test.expectedRes, res)
		})
	}
}
