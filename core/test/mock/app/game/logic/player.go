// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/game/logic/contract/player.go

// Package mock_game_logic is a generated GoMock package.
package mock_game_logic

import (
	reflect "reflect"
	contract "uwwolf/internal/app/game/logic/contract"
	types "uwwolf/internal/app/game/logic/types"

	gomock "github.com/golang/mock/gomock"
	orb "github.com/paulmach/orb"
)

// MockPlayer is a mock of Player interface.
type MockPlayer struct {
	ctrl     *gomock.Controller
	recorder *MockPlayerMockRecorder
}

// MockPlayerMockRecorder is the mock recorder for MockPlayer.
type MockPlayerMockRecorder struct {
	mock *MockPlayer
}

// NewMockPlayer creates a new mock instance.
func NewMockPlayer(ctrl *gomock.Controller) *MockPlayer {
	mock := &MockPlayer{ctrl: ctrl}
	mock.recorder = &MockPlayerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlayer) EXPECT() *MockPlayerMockRecorder {
	return m.recorder
}

// AssignRole mocks base method.
func (m *MockPlayer) AssignRole(roleId types.RoleId) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignRole", roleId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignRole indicates an expected call of AssignRole.
func (mr *MockPlayerMockRecorder) AssignRole(roleId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignRole", reflect.TypeOf((*MockPlayer)(nil).AssignRole), roleId)
}

// Die mocks base method.
func (m *MockPlayer) Die() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Die")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Die indicates an expected call of Die.
func (mr *MockPlayerMockRecorder) Die() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Die", reflect.TypeOf((*MockPlayer)(nil).Die))
}

// Exit mocks base method.
func (m *MockPlayer) Exit() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exit")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Exit indicates an expected call of Exit.
func (mr *MockPlayerMockRecorder) Exit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exit", reflect.TypeOf((*MockPlayer)(nil).Exit))
}

// FactionId mocks base method.
func (m *MockPlayer) FactionId() types.FactionId {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FactionId")
	ret0, _ := ret[0].(types.FactionId)
	return ret0
}

// FactionId indicates an expected call of FactionId.
func (mr *MockPlayerMockRecorder) FactionId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FactionId", reflect.TypeOf((*MockPlayer)(nil).FactionId))
}

// Id mocks base method.
func (m *MockPlayer) Id() types.PlayerId {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(types.PlayerId)
	return ret0
}

// Id indicates an expected call of Id.
func (mr *MockPlayerMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockPlayer)(nil).Id))
}

// IsDead mocks base method.
func (m *MockPlayer) IsDead() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDead")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDead indicates an expected call of IsDead.
func (mr *MockPlayerMockRecorder) IsDead() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDead", reflect.TypeOf((*MockPlayer)(nil).IsDead))
}

// Location mocks base method.
func (m *MockPlayer) Location() (float64, float64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Location")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(float64)
	return ret0, ret1
}

// Location indicates an expected call of Location.
func (mr *MockPlayerMockRecorder) Location() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Location", reflect.TypeOf((*MockPlayer)(nil).Location))
}

// MainRoleId mocks base method.
func (m *MockPlayer) MainRoleId() types.RoleId {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MainRoleId")
	ret0, _ := ret[0].(types.RoleId)
	return ret0
}

// MainRoleId indicates an expected call of MainRoleId.
func (mr *MockPlayerMockRecorder) MainRoleId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MainRoleId", reflect.TypeOf((*MockPlayer)(nil).MainRoleId))
}

// Move mocks base method.
func (m *MockPlayer) Move(position orb.Point) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Move", position)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Move indicates an expected call of Move.
func (mr *MockPlayerMockRecorder) Move(position interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Move", reflect.TypeOf((*MockPlayer)(nil).Move), position)
}

// PlayRecords mocks base method.
func (m *MockPlayer) PlayRecords() []types.PlayerRecord {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlayRecords")
	ret0, _ := ret[0].([]types.PlayerRecord)
	return ret0
}

// PlayRecords indicates an expected call of PlayRecords.
func (mr *MockPlayerMockRecorder) PlayRecords() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlayRecords", reflect.TypeOf((*MockPlayer)(nil).PlayRecords))
}

// RevokeRole mocks base method.
func (m *MockPlayer) RevokeRole(roleId types.RoleId) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeRole", roleId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RevokeRole indicates an expected call of RevokeRole.
func (mr *MockPlayerMockRecorder) RevokeRole(roleId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeRole", reflect.TypeOf((*MockPlayer)(nil).RevokeRole), roleId)
}

// RoleIds mocks base method.
func (m *MockPlayer) RoleIds() []types.RoleId {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoleIds")
	ret0, _ := ret[0].([]types.RoleId)
	return ret0
}

// RoleIds indicates an expected call of RoleIds.
func (mr *MockPlayerMockRecorder) RoleIds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoleIds", reflect.TypeOf((*MockPlayer)(nil).RoleIds))
}

// Roles mocks base method.
func (m *MockPlayer) Roles() map[types.RoleId]contract.Role {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Roles")
	ret0, _ := ret[0].(map[types.RoleId]contract.Role)
	return ret0
}

// Roles indicates an expected call of Roles.
func (mr *MockPlayerMockRecorder) Roles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Roles", reflect.TypeOf((*MockPlayer)(nil).Roles))
}

// SetFactionId mocks base method.
func (m *MockPlayer) SetFactionId(factionId types.FactionId) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFactionId", factionId)
}

// SetFactionId indicates an expected call of SetFactionId.
func (mr *MockPlayerMockRecorder) SetFactionId(factionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFactionId", reflect.TypeOf((*MockPlayer)(nil).SetFactionId), factionId)
}

// UseRole mocks base method.
func (m *MockPlayer) UseRole(req types.RoleRequest) types.RoleResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UseRole", req)
	ret0, _ := ret[0].(types.RoleResponse)
	return ret0
}

// UseRole indicates an expected call of UseRole.
func (mr *MockPlayerMockRecorder) UseRole(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UseRole", reflect.TypeOf((*MockPlayer)(nil).UseRole), req)
}
