package action

import (
	"fmt"
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

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
		expectedErr string
		setup       func(*gamemock.MockGame, *gamemock.MockPoll)
	}{
		{
			name: "Failure (Poll doesnt exist)",
			setting: &VoteActionSetting{
				FactionID: vs.factionID,
			},
			expectedErr: "Poll does not exist ¯\\_(ツ)_/¯",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vs.factionID).Return(nil).Times(1)
			},
		},
		{
			name: "Ok",
			setting: &VoteActionSetting{
				FactionID: vs.factionID,
				PlayerID:  vs.actorID,
				Weight:    vs.weight,
			},
			expectedErr: "",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mp.EXPECT().AddElectors(vs.actorID).Times(1)
				mp.EXPECT().SetWeight(vs.actorID, vs.weight).Return(true).Times(1)
				mg.EXPECT().Poll(vs.factionID).Return(mp).Times(2)
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			game := gamemock.NewMockGame(ctrl)
			poll := gamemock.NewMockPoll(ctrl)
			test.setup(game, poll)

			v, err := NewVote(game, test.setting)

			if test.expectedErr != "" {
				vs.Nil(v)
				vs.Equal(test.expectedErr, err.Error())
			} else {
				vs.Equal(vars.VoteActionID, v.ID())
				vs.NotNil(v.(*vote).poll)
				vs.Equal(poll, v.(*vote).poll)
			}
		})
	}
}

func (vs VoteSuite) TestSkipVote() {
	ctrl := gomock.NewController(vs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	poll := gamemock.NewMockPoll(ctrl)

	game.EXPECT().Poll(vs.factionID).Return(poll).Times(2)
	poll.EXPECT().AddElectors(vs.actorID).Times(1)
	poll.EXPECT().SetWeight(vs.actorID, vs.weight).Return(true).Times(1)
	poll.EXPECT().Vote(vs.actorID, types.PlayerID("")).Times(1)

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

func (vs VoteSuite) TestPerformVote() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
		setup       func(*gamemock.MockPoll)
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
			setup: func(mp *gamemock.MockPoll) {
				mp.EXPECT().Vote(vs.actorID, vs.targetID).
					Return(false, fmt.Errorf("CANT_VOTE error")).Times(1)
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
			setup: func(mp *gamemock.MockPoll) {
				mp.EXPECT().Vote(vs.actorID, vs.targetID).Return(true, nil).Times(1)
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			game := gamemock.NewMockGame(ctrl)
			poll := gamemock.NewMockPoll(ctrl)

			game.EXPECT().Poll(vs.factionID).Return(poll).Times(2)
			poll.EXPECT().AddElectors(vs.actorID).Times(1)
			poll.EXPECT().SetWeight(vs.actorID, vs.weight).Return(true).Times(1)
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
