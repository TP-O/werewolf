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

// type RecognizeSuite struct {
// 	suite.Suite
// 	ctrl          *gomock.Controller
// 	game          *gamemock.MockGame
// 	actorID       enum.PlayerID
// 	recognizedIDs []enum.PlayerID
// 	factionID     enum.FactionID
// 	roleID        enum.RoleID
// }

// func TestRecognizeSuite(t *testing.T) {
// 	suite.Run(t, new(RecognizeSuite))
// }

// func (rs *RecognizeSuite) SetupSuite() {
// 	rs.actorID = "1"
// 	rs.recognizedIDs = []enum.PlayerID{"1", "2"}
// 	rs.factionID = enum.WerewolfFactionID
// 	rs.roleID = enum.WerewolfRoleID
// }

// func (rs *RecognizeSuite) SetupTest() {
// 	rs.ctrl = gomock.NewController(rs.T())
// 	rs.game = gamemock.NewMockGame(rs.ctrl)
// }

// func (rs *RecognizeSuite) TearDownTest() {
// 	rs.ctrl.Finish()
// }

// func (rs *RecognizeSuite) TestNewFactionRecognize() {
// 	factionRecognize := NewFactionRecognize(rs.game, rs.factionID)

// 	rs.Equal(enum.RecognizeActionID, factionRecognize.ID())
// 	rs.NotNil(factionRecognize.State())
// 	rs.IsType(new(types.RecognizeState), factionRecognize.State())
// 	rs.NotNil(factionRecognize.State().(*types.RecognizeState).Faction)
// 	rs.Nil(factionRecognize.State().(*types.RecognizeState).Role)
// }

// func (rs *RecognizeSuite) TestNewRoleRecognize() {
// 	roleRecognize := NewRoleRecognize(rs.game, rs.roleID)

// 	rs.Equal(enum.RecognizeActionID, roleRecognize.ID())
// 	rs.NotNil(roleRecognize.State())
// 	rs.IsType(new(types.RecognizeState), roleRecognize.State())
// 	rs.NotNil(roleRecognize.State().(*types.RecognizeState).Role)
// 	rs.Nil(roleRecognize.State().(*types.RecognizeState).Faction)
// }

// func (rs *RecognizeSuite) TestPerformFactionRecognize() {
// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedRes *types.ActionResponse
// 		newState    []enum.PlayerID
// 		setup       func(contract.Action)
// 	}{
// 		{
// 			name: "Ok",
// 			req: &types.ActionRequest{
// 				ActorID: rs.actorID,
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      rs.recognizedIDs,
// 				Message:   "",
// 			},
// 			newState: rs.recognizedIDs,
// 			setup: func(contract.Action) {
// 				rs.game.EXPECT().PlayerIDsByFactionID(rs.factionID).Return(rs.recognizedIDs).Times(1)
// 			},
// 		},
// 		{
// 			name: "Ok (Return cache in the second time)",
// 			req: &types.ActionRequest{
// 				ActorID: rs.actorID,
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      rs.recognizedIDs,
// 				Message:   "",
// 			},
// 			newState: rs.recognizedIDs,
// 			setup: func(recognizeAction contract.Action) {
// 				recognizeAction.(*recognize).isRecognized = true
// 				recognizeAction.(*recognize).state.Faction = rs.recognizedIDs
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		rs.Run(test.name, func() {
// 			recognize := NewFactionRecognize(rs.game, rs.factionID)
// 			test.setup(recognize)
// 			res := recognize.Execute(test.req)

// 			rs.Equal(test.expectedRes, res)
// 			rs.Equal(test.newState, recognize.State().(*types.RecognizeState).Faction)
// 		})
// 	}
// }

// func (rs *RecognizeSuite) TestPerformRoleRecognize() {
// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedRes *types.ActionResponse
// 		newState    []enum.PlayerID
// 		setup       func(contract.Action)
// 	}{
// 		{
// 			name: "Ok",
// 			req: &types.ActionRequest{
// 				ActorID: rs.actorID,
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      rs.recognizedIDs,
// 				Message:   "",
// 			},
// 			newState: rs.recognizedIDs,
// 			setup: func(contract.Action) {
// 				rs.game.EXPECT().PlayerIDsByRoleID(rs.roleID).Return(rs.recognizedIDs).Times(1)
// 			},
// 		},
// 		{
// 			name: "Ok but second time (return cache)",
// 			req: &types.ActionRequest{
// 				ActorID: rs.actorID,
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: false,
// 				Data:      rs.recognizedIDs,
// 				Message:   "",
// 			},
// 			newState: rs.recognizedIDs,
// 			setup: func(recognizeAction contract.Action) {
// 				recognizeAction.(*recognize).isRecognized = true
// 				recognizeAction.(*recognize).state.Role = rs.recognizedIDs
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		rs.Run(test.name, func() {
// 			recognize := NewRoleRecognize(rs.game, rs.roleID)
// 			test.setup(recognize)
// 			res := recognize.Execute(test.req)

// 			rs.Equal(test.expectedRes, res)
// 			rs.Equal(test.newState, recognize.State().(*types.RecognizeState).Role)
// 		})
// 	}
// }
