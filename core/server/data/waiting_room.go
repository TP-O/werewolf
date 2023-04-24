package data

import (
	"uwwolf/game/types"
)

type WaitingRoom struct {
	ID        string           `json:"id"`
	OwnerID   types.PlayerID   `json:"owner_id"`
	PlayerIDs []types.PlayerID `json:"player_ids"`
}
