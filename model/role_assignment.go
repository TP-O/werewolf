package model

import "uwwolf/game/enum"

type RoleAssignment struct {
	GameID   enum.GameID   `cql:"game_id" json:"gameID"`
	PlayerID enum.PlayerID `cql:"player_id" json:"playerID"`
	RoleID   enum.RoleID   `cql:"role_id" json:"roleID"`
	IsLeader bool          `cql:"is_leader" json:"isLeader"`
	IsDead   bool          `cql:"is_dead" json:"isDead"`
}
