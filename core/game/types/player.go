package types

// PlayerID is ID type of player.
type PlayerID string

// IsUnknown checks if player ID is empty.
func (p PlayerID) IsUnknown() bool {
	return p == ""
}
