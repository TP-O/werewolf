package action

// import (
// 	"testing"
// 	"uwwolf/game/contract"
// 	"uwwolf/game/enum"
// 	"uwwolf/game/types"
// 	gamemock "uwwolf/mock/game"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type PredictSuite struct {
// 	suite.Suite
// 	ctrl           *gomock.Controller
// 	game           *gamemock.MockGame
// 	targetedPlayer *gamemock.MockPlayer
// 	actorID        enum.PlayerID
// 	targetID       enum.PlayerID
// 	factionID      enum.FactionID
// 	roleID         enum.RoleID
// }

// func TestPredictSuite(t *testing.T) {
// 	suite.Run(t, new(PredictSuite))
// }

// func (ps *PredictSuite) SetupSuite() {
// 	ps.actorID = "1"
// 	ps.targetID = "2"
// 	ps.factionID = enum.WerewolfFactionID
// 	ps.roleID = enum.WerewolfRoleID
// }

// func (ps *PredictSuite) SetupTest() {
// 	ps.ctrl = gomock.NewController(ps.T())
// 	ps.game = gamemock.NewMockGame(ps.ctrl)
// 	ps.targetedPlayer = gamemock.NewMockPlayer(ps.ctrl)
// }

// func (ps *PredictSuite) TearDownTest() {
// 	ps.ctrl.Finish()
// }

// func (ps *PredictSuite) TestNewFactionPredict() {
// 	factionPredict := NewFactionPredict(ps.game, ps.factionID)

// 	ps.Equal(enum.PredictActionID, factionPredict.ID())
// 	ps.NotNil(factionPredict.State())
// 	ps.IsType(new(types.PredictState), factionPredict.State())
// 	ps.NotNil(factionPredict.State().(*types.PredictState).Faction)
// 	ps.Nil(factionPredict.State().(*types.PredictState).Role)

// }

// func (ps *PredictSuite) TestNewRolePredict() {
// 	rolePredict := NewRolePredict(ps.game, ps.roleID)

// 	ps.Equal(enum.PredictActionID, rolePredict.ID())
// 	ps.NotNil(rolePredict.State())
// 	ps.IsType(new(types.PredictState), rolePredict.State())
// 	ps.NotNil(rolePredict.State().(*types.PredictState).Role)
// 	ps.Nil(rolePredict.State().(*types.PredictState).Faction)

// }

// func (ps *PredictSuite) TestValidateFactionPredict() {
// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedErr string
// 		setup       func(contract.Action)
// 	}{
// 		{
// 			name:        "Failure (Empty action request)",
// 			req:         nil,
// 			expectedErr: "Action request can not be empty (⊙＿⊙')",
// 			setup:       func(contract.Action) {},
// 		},
// 		{
// 			name: "Failure (Cannot predict myself)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{ps.actorID},
// 			},
// 			expectedErr: "WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻",
// 			setup:       func(contract.Action) {},
// 		},
// 		{
// 			name: "Failure (Cannot predict known player)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{ps.targetID},
// 			},
// 			expectedErr: "You already knew this player ¯\\(º_o)/¯",
// 			setup: func(predict contract.Action) {
// 				predict.State().(*types.PredictState).Faction[ps.targetID] = ps.factionID
// 			},
// 		},
// 		{
// 			name: "Failure (Cannot predict non-existent player)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{"-99"},
// 			},
// 			expectedErr: "Non-existent player ¯\\_(ツ)_/¯",
// 			setup: func(predict contract.Action) {
// 				ps.game.EXPECT().Player("-99").Return(nil).Times(1)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			factionPredict := NewFactionPredict(ps.game, ps.factionID)
// 			test.setup(factionPredict)
// 			err := factionPredict.Validate(test.req)

// 			ps.Equal(test.expectedErr, err.Error())
// 		})
// 	}
// }

// func (ps *PredictSuite) TestValidateRolePredict() {
// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedErr string
// 		setup       func(contract.Action)
// 	}{
// 		{
// 			name:        "Failure (Empty action request)",
// 			req:         nil,
// 			expectedErr: "Action request can not be empty (⊙＿⊙')",
// 			setup:       func(contract.Action) {},
// 		},
// 		{
// 			name: "Failure (Cannot predict myself)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{ps.actorID},
// 			},
// 			expectedErr: "WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻",
// 			setup:       func(contract.Action) {},
// 		},
// 		{
// 			name: "Failure (Cannot predict known player)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{ps.targetID},
// 			},
// 			expectedErr: "You already knew this player ¯\\(º_o)/¯",
// 			setup: func(predict contract.Action) {
// 				predict.State().(*types.PredictState).Role[ps.targetID] = ps.roleID
// 			},
// 		},
// 		{
// 			name: "Failure (Cannot predict non-existent player)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{"-99"},
// 			},
// 			expectedErr: "Non-existent player ¯\\_(ツ)_/¯",
// 			setup: func(predict contract.Action) {
// 				ps.game.EXPECT().Player("-99").Return(nil).Times(1)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			rolePredict := NewRolePredict(ps.game, ps.roleID)
// 			test.setup(rolePredict)
// 			err := rolePredict.Validate(test.req)

// 			ps.Equal(test.expectedErr, err.Error())
// 		})
// 	}
// }

// func (ps *PredictSuite) TestPerformFactionPredict() {
// 	ps.game.EXPECT().Player(ps.targetID).Return(ps.targetedPlayer).AnyTimes()

// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedRes *types.ActionResponse
// 		newState    map[enum.PlayerID]enum.FactionID
// 		setup       func(contract.Action)
// 	}{
// 		{
// 			name: "Ok (Incorrect prediction)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{ps.targetID},
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      false,
// 				Message:   "",
// 			},
// 			newState: map[enum.PlayerID]enum.FactionID{
// 				ps.targetID: enum.FactionID(0),
// 			},
// 			setup: func(contract.Action) {
// 				ps.targetedPlayer.EXPECT().FactionID().Return(enum.VillagerFactionID)
// 				ps.targetedPlayer.EXPECT().ID().Return(ps.targetID)
// 			},
// 		},
// 		{
// 			name: "Ok (Correct prediction)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{ps.targetID},
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      true,
// 				Message:   "",
// 			},
// 			newState: map[enum.PlayerID]enum.FactionID{
// 				ps.targetID: ps.factionID,
// 			},
// 			setup: func(contract.Action) {
// 				ps.targetedPlayer.EXPECT().FactionID().Return(enum.WerewolfFactionID)
// 				ps.targetedPlayer.EXPECT().ID().Return(ps.targetID)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			factionPredict := NewFactionPredict(ps.game, ps.factionID)
// 			test.setup(factionPredict)
// 			res := factionPredict.Perform(test.req)

// 			ps.Equal(test.expectedRes, res)
// 			ps.Equal(test.newState, factionPredict.State().(*types.PredictState).Faction)
// 		})
// 	}
// }

// func (ps *PredictSuite) TestPerformRolePredict() {
// 	ps.game.EXPECT().Player(ps.targetID).Return(ps.targetedPlayer).AnyTimes()

// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedRes *types.ActionResponse
// 		newState    map[enum.PlayerID]enum.RoleID
// 		setup       func(contract.Action)
// 	}{
// 		{
// 			name: "Ok (Incorrect prediction)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{ps.targetID},
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      false,
// 				Message:   "",
// 			},
// 			newState: map[enum.PlayerID]enum.RoleID{
// 				ps.targetID: enum.RoleID(0),
// 			},
// 			setup: func(contract.Action) {
// 				ps.targetedPlayer.EXPECT().RoleIDs().Return([]enum.RoleID{
// 					enum.VillagerRoleID,
// 				})
// 				ps.targetedPlayer.EXPECT().ID().Return(ps.targetID)
// 			},
// 		},
// 		{
// 			name: "Ok (Correct prediction)",
// 			req: &types.ActionRequest{
// 				ActorID:   ps.actorID,
// 				TargetIDs: []enum.PlayerID{ps.targetID},
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      true,
// 				Message:   "",
// 			},
// 			newState: map[enum.PlayerID]enum.RoleID{
// 				ps.targetID: ps.roleID,
// 			},
// 			setup: func(contract.Action) {
// 				ps.targetedPlayer.EXPECT().RoleIDs().Return([]enum.RoleID{
// 					enum.WerewolfRoleID,
// 				})
// 				ps.targetedPlayer.EXPECT().ID().Return(ps.targetID)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			rolePredict := NewRolePredict(ps.game, ps.roleID)
// 			test.setup(rolePredict)
// 			res := rolePredict.Perform(test.req)

// 			ps.Equal(test.expectedRes, res)
// 			ps.Equal(test.newState, rolePredict.State().(*types.PredictState).Role)
// 		})
// 	}
// }
