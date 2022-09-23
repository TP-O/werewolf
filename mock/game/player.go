// Code generated by MockGen. DO NOT EDIT.
// Source: app/game/contract/player.go

// Package game is a generated GoMock package.
package game

import (
	reflect "reflect"
	contract "uwwolf/app/game/contract"
	types "uwwolf/app/types"

	gomock "github.com/golang/mock/gomock"
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

// AssignRoles mocks base method.
func (m *MockPlayer) AssignRoles(roles ...contract.Role) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range roles {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "AssignRoles", varargs...)
}

// AssignRoles indicates an expected call of AssignRoles.
func (mr *MockPlayerMockRecorder) AssignRoles(roles ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignRoles", reflect.TypeOf((*MockPlayer)(nil).AssignRoles), roles...)
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

// UseSkill mocks base method.
func (m *MockPlayer) UseSkill(req *types.ActionRequest) *types.ActionResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UseSkill", req)
	ret0, _ := ret[0].(*types.ActionResponse)
	return ret0
}

// UseSkill indicates an expected call of UseSkill.
func (mr *MockPlayerMockRecorder) UseSkill(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UseSkill", reflect.TypeOf((*MockPlayer)(nil).UseSkill), req)
}
