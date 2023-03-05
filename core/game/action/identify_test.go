package action

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RecognizeSuite struct {
	suite.Suite
	identifiedFactionID types.FactionID
	identifiedRoleID    types.RoleID
	actorID             types.PlayerID
	identifiedIDs       []types.PlayerID
}

func TestRecognizeSuite(t *testing.T) {
	suite.Run(t, new(RecognizeSuite))
}

func (rs *RecognizeSuite) SetupSuite() {
	rs.actorID = types.PlayerID("1")
	rs.identifiedIDs = []types.PlayerID{"1", "2"}
	rs.identifiedFactionID = vars.WerewolfFactionID
	rs.identifiedRoleID = vars.WerewolfRoleID
}

func (rs *RecognizeSuite) TestNewFactionIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)

	fIdent := NewFactionIdentify(game, rs.identifiedFactionID).(*identify)

	rs.Equal(vars.IdentifyActionID, fIdent.ID())
	rs.Equal(rs.identifiedFactionID, fIdent.FactionID)
	rs.Len(fIdent.Faction, 0)
	rs.Equal(types.RoleID(0), fIdent.RoleID)
	rs.Len(fIdent.Role, 0)
	rs.False(fIdent.isIdentified)
}

func (rs *RecognizeSuite) TestNewRoleIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)

	ident := NewRoleIdentify(game, rs.identifiedRoleID).(*identify)

	rs.Equal(vars.IdentifyActionID, ident.ID())
	rs.Equal(types.FactionID(0), ident.FactionID)
	rs.Len(ident.Faction, 0)
	rs.Equal(rs.identifiedRoleID, ident.RoleID)
	rs.Len(ident.Role, 0)
	rs.False(ident.isIdentified)
}

func (rs *RecognizeSuite) TestValidateFactionIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)

	expectErr := "You already recognized everyone ¯\\(º_o)/¯"

	ident := NewFactionIdentify(game, rs.identifiedFactionID).(*identify)
	ident.isIdentified = true

	err := ident.validate(&types.ActionRequest{})
	rs.Equal(expectErr, err.Error())
}

func (rs *RecognizeSuite) TestPerformFactionIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)

	game.EXPECT().PlayerIDsWithFactionID(rs.identifiedFactionID, false).
		Return(rs.identifiedIDs).Times(1)

	req := &types.ActionRequest{
		ActorID: rs.actorID,
	}
	expectRes := &types.ActionResponse{
		Ok:   true,
		Data: rs.identifiedIDs,
	}

	ident := NewFactionIdentify(game, rs.identifiedFactionID).(*identify)
	res := ident.perform(req)

	rs.Equal(expectRes, res)
	rs.Equal(rs.identifiedIDs, ident.Faction)
	rs.True(ident.isIdentified)
}

func (rs *RecognizeSuite) TestValidateRoleIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)

	expectErr := "You already recognized everyone ¯\\(º_o)/¯"

	ident := NewRoleIdentify(game, rs.identifiedRoleID).(*identify)
	ident.isIdentified = true

	err := ident.validate(&types.ActionRequest{})
	rs.Equal(expectErr, err.Error())
}

func (rs *RecognizeSuite) TestPerformRoleIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)

	game.EXPECT().PlayerIDsWithRoleID(rs.identifiedRoleID).
		Return(rs.identifiedIDs).Times(1)

	req := &types.ActionRequest{
		ActorID: rs.actorID,
	}
	expectRes := &types.ActionResponse{
		Ok:   true,
		Data: rs.identifiedIDs,
	}

	ident := NewRoleIdentify(game, rs.identifiedRoleID).(*identify)
	res := ident.perform(req)

	rs.Equal(expectRes, res)
	rs.Equal(rs.identifiedIDs, ident.Role)
	rs.True(ident.isIdentified)
}
