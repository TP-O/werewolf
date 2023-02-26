package action

// import (
// 	"fmt"
// 	"testing"
// 	"uwwolf/game/enum"
// 	"uwwolf/game/types"
// 	gamemock "uwwolf/mock/game"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type VoteSuite struct {
// 	suite.Suite
// 	ctrl      *gomock.Controller
// 	game      *gamemock.MockGame
// 	poll      *gamemock.MockPoll
// 	actorID   enum.PlayerID
// 	targetID  enum.PlayerID
// 	factionID enum.FactionID
// }

// func TestVoteSuite(t *testing.T) {
// 	suite.Run(t, new(VoteSuite))
// }

// func (vs *VoteSuite) SetupSuite() {
// 	vs.actorID = "2"
// 	vs.targetID = "2"
// }

// func (vs *VoteSuite) SetupTest() {
// 	vs.ctrl = gomock.NewController(vs.T())
// 	vs.game = gamemock.NewMockGame(vs.ctrl)
// 	vs.poll = gamemock.NewMockPoll(vs.ctrl)
// }

// func (vs *VoteSuite) TearDownTest() {
// 	vs.ctrl.Finish()
// }

// func (vs *VoteSuite) TestNewVote() {
// 	weight := uint(2)
// 	tests := []struct {
// 		name        string
// 		setting     *types.VoteActionSetting
// 		expectedErr string
// 		setup       func()
// 	}{
// 		{
// 			name: "Failure (Poll does not exist)",
// 			setting: &types.VoteActionSetting{
// 				FactionID: vs.factionID,
// 				PlayerID:  vs.actorID,
// 				Weight:    weight,
// 			},
// 			expectedErr: "Poll does not exist ¯\\_(ツ)_/¯",
// 			setup: func() {
// 				vs.game.EXPECT().Poll(vs.factionID).Return(nil).Times(1)
// 			},
// 		},
// 		{
// 			name: "Failure (Cannot add player as elector to poll)",
// 			setting: &types.VoteActionSetting{
// 				FactionID: vs.factionID,
// 				PlayerID:  vs.actorID,
// 				Weight:    weight,
// 			},
// 			expectedErr: "Unable to join to the poll ಠ_ಠ",
// 			setup: func() {
// 				vs.game.EXPECT().Poll(vs.factionID).Return(vs.poll).Times(1)
// 				vs.poll.EXPECT().AddElectors(vs.actorID).Return(false).Times(1)
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			setting: &types.VoteActionSetting{
// 				FactionID: vs.factionID,
// 				PlayerID:  vs.actorID,
// 				Weight:    weight,
// 			},
// 			expectedErr: "",
// 			setup: func() {
// 				vs.game.EXPECT().Poll(vs.factionID).Return(vs.poll).Times(1)
// 				vs.poll.EXPECT().AddElectors(vs.actorID).Return(true).Times(1)
// 				vs.poll.EXPECT().SetWeight(vs.actorID, weight).Return(true).Times(1)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		vs.Run(test.name, func() {
// 			test.setup()
// 			vote, err := NewVote(vs.game, test.setting)

// 			if test.expectedErr != "" {
// 				vs.Nil(vote)
// 				vs.Equal(test.expectedErr, err.Error())
// 			} else {
// 				vs.Equal(enum.VoteActionID, vote.ID())
// 				vs.NotNil(vote.State())
// 			}
// 		})
// 	}
// }

// func (vs *VoteSuite) TestSkipVote() {
// 	// Mock for vote action initialization
// 	weight := uint(2)
// 	vs.game.EXPECT().Poll(vs.factionID).Return(vs.poll).Times(1)
// 	vs.poll.EXPECT().AddElectors(vs.actorID).Return(true).Times(1)
// 	vs.poll.EXPECT().SetWeight(vs.actorID, weight).Return(true).Times(1)
// 	vs.poll.EXPECT().Vote(vs.actorID, enum.PlayerID("")).Times(1)

// 	expectedRes := &types.ActionResponse{
// 		Ok:        true,
// 		IsSkipped: true,
// 		Data:      nil,
// 		Message:   "",
// 	}
// 	vote, _ := NewVote(vs.game, &types.VoteActionSetting{
// 		FactionID: vs.factionID,
// 		PlayerID:  vs.actorID,
// 		Weight:    weight,
// 	})
// 	res := vote.Skip(&types.ActionRequest{
// 		ActorID:   vs.actorID,
// 		IsSkipped: true,
// 	})

// 	vs.Equal(expectedRes, res)
// }

// func (vs *VoteSuite) TestPerformVote() {
// 	weight := uint(2)
// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedRes *types.ActionResponse
// 		setup       func()
// 	}{

// 		{
// 			name: "Vote failed",
// 			req: &types.ActionRequest{
// 				ActorID:   vs.actorID,
// 				TargetIDs: []enum.PlayerID{vs.targetID},
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        false,
// 				IsSkipped: false,
// 				Data:      nil,
// 				Message:   "CANT_VOTE error",
// 			},
// 			setup: func() {
// 				vs.poll.
// 					EXPECT().
// 					Vote(vs.actorID, vs.targetID).
// 					Return(false, fmt.Errorf("CANT_VOTE error")).
// 					Times(1)
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			req: &types.ActionRequest{
// 				ActorID:   vs.actorID,
// 				TargetIDs: []enum.PlayerID{vs.targetID},
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      vs.targetID,
// 				Message:   "",
// 			},
// 			setup: func() {
// 				vs.poll.
// 					EXPECT().
// 					Vote(vs.actorID, vs.targetID).
// 					Return(true, nil).
// 					Times(1)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		vs.Run(test.name, func() {
// 			// Mock for vote action initialization
// 			vs.game.EXPECT().Poll(vs.factionID).Return(vs.poll).Times(1)
// 			vs.poll.EXPECT().AddElectors(vs.actorID).Return(true).Times(1)
// 			vs.poll.EXPECT().SetWeight(vs.actorID, weight).Return(true).Times(1)

// 			vote, _ := NewVote(vs.game, &types.VoteActionSetting{
// 				FactionID: vs.factionID,
// 				PlayerID:  vs.actorID,
// 				Weight:    weight,
// 			})
// 			test.setup()
// 			res := vote.Perform(test.req)

// 			vs.Equal(test.expectedRes, res)
// 		})
// 	}
// }
