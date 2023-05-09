package model

import (
	"uwwolf/internal/app/game/logic/types"
)

type WaitingRoom struct {
	ID        string           `json:"id"`
	OwnerID   types.PlayerID   `json:"owner_id"`
	PlayerIDs []types.PlayerID `json:"player_ids"`
}
