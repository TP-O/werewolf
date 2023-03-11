package action

import (
	"fmt"
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	mock_game "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type VoteSuite struct {
	suite.Suite
	actorID   types.PlayerID
	targetID  types.PlayerID
	factionID types.FactionID
	weight    uint
}

func TestVoteSuite(t *testing.T) {
	suite.Run(t, new(VoteSuite))
}

func (vs *VoteSuite) SetupSuite() {
	vs.actorID = types.PlayerID("2")
	vs.targetID = types.PlayerID("2")
	vs.factionID = vars.VillagerFactionID
	vs.weight = 1
}

func (vs VoteSuite) TestNewVote() {
	tests := []struct {
		name        string
		setting     *VoteActionSetting
		expectedErr error
		setup       func(*mock_game.MockGame, *mock_game.MockPoll)
	}{
		{
			name: "Failure (Poll doesnt exist)",
			setting: &VoteActionSetting{
				FactionID: vs.factionID,
			},
			expectedErr: fmt.Errorf("Poll does not exist ¯\\_(ツ)_/¯"),
			setup: func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {
				mg.EXPECT().Poll(vs.factionID).Return(nil)
			},
		},
		{
			name: "Ok",
			setting: &VoteActionSetting{
				FactionID: vs.factionID,
				PlayerID:  vs.actorID,
				Weight:    vs.weight,
			},
			setup: func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {
				mp.EXPECT().AddElectors(vs.actorID)
				mp.EXPECT().SetWeight(vs.actorID, vs.weight).Return(true)
				mg.EXPECT().Poll(vs.factionID).Return(mp).Times(2)
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			game := mock_game.NewMockGame(ctrl)
			poll := mock_game.NewMockPoll(ctrl)
			test.setup(game, poll)

			v, err := NewVote(game, test.setting)

			if test.expectedErr != nil {
				vs.Nil(v)
				vs.Equal(test.expectedErr, err)
			} else {
				vs.Equal(vars.VoteActionID, v.ID())
				vs.NotNil(v.(*vote).poll)
				vs.Equal(poll, v.(*vote).poll)
			}
		})
	}
}

func (vs VoteSuite) TestSkip() {
	ctrl := gomock.NewController(vs.T())
	defer ctrl.Finish()
	game := mock_game.NewMockGame(ctrl)
	poll := mock_game.NewMockPoll(ctrl)

	game.EXPECT().Poll(vs.factionID).Return(poll).Times(2)
	poll.EXPECT().AddElectors(vs.actorID)
	poll.EXPECT().SetWeight(vs.actorID, vs.weight).Return(true)
	poll.EXPECT().Vote(vs.actorID, types.PlayerID(""))

	expectedRes := &types.ActionResponse{
		Ok:        true,
		IsSkipped: true,
		Message:   "Skipped!",
	}

	v, _ := NewVote(game, &VoteActionSetting{
		FactionID: vs.factionID,
		PlayerID:  vs.actorID,
		Weight:    vs.weight,
	})
	res := v.(*vote).skip(&types.ActionRequest{
		ActorID:   vs.actorID,
		IsSkipped: true,
	})

	vs.Equal(expectedRes, res)
}

func (vs VoteSuite) TestPerform() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
		setup       func(*mock_game.MockPoll)
	}{

		{
			name: "Failure (Vote failed)",
			req: &types.ActionRequest{
				ActorID:  vs.actorID,
				TargetID: vs.targetID,
			},
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "CANT_VOTE error",
			},
			setup: func(mp *mock_game.MockPoll) {
				mp.EXPECT().Vote(vs.actorID, vs.targetID).
					Return(false, fmt.Errorf("CANT_VOTE error"))
			},
		},
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorID:  vs.actorID,
				TargetID: vs.targetID,
			},
			expectedRes: &types.ActionResponse{
				Ok:   true,
				Data: vs.targetID,
			},
			setup: func(mp *mock_game.MockPoll) {
				mp.EXPECT().Vote(vs.actorID, vs.targetID).Return(true, nil)
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			game := mock_game.NewMockGame(ctrl)
			poll := mock_game.NewMockPoll(ctrl)

			game.EXPECT().Poll(vs.factionID).Return(poll).Times(2)
			poll.EXPECT().AddElectors(vs.actorID)
			poll.EXPECT().SetWeight(vs.actorID, vs.weight).Return(true)
			test.setup(poll)

			v, _ := NewVote(game, &VoteActionSetting{
				FactionID: vs.factionID,
				PlayerID:  vs.actorID,
				Weight:    vs.weight,
			})
			res := v.(*vote).perform(test.req)

			vs.Equal(test.expectedRes, res)
		})
	}
}
