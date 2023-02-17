package model

import "uwwolf/game/enum"

type Player struct {
	ID       enum.PlayerID     `cql:"id" json:"id"`
	StatusID enum.PlayerStatus `cql:"status_id" json:"status_id"`
}
