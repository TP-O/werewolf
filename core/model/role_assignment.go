package model

import "uwwolf/game/types"

type RoleAssignment struct {
	GameID   types.GameID   `cql:"game_id" json:"gameID"`
	PlayerID types.PlayerID `cql:"player_id" json:"playerID"`
	RoleID   types.RoleID   `cql:"role_id" json:"roleID"`
	IsLeader bool           `cql:"is_leader" json:"isLeader"`
	IsDead   bool           `cql:"is_dead" json:"isDead"`
}
