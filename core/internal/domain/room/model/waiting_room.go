package model

import (
	"uwwolf/internal/app/game/logic/types"
)

type WaitingRoom struct {
	ID        string           `json:"id"`
	OwnerID   types.PlayerId   `json:"owner_id"`
	PlayerIDs []types.PlayerId `json:"player_ids"`
}
