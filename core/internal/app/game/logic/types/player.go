package types

// PlayerID is ID type of player.
type PlayerId = string

// PlayerRecord is the player's play record
type PlayerRecord struct {
	Round
	Turn
	RoleId
	ActionId
	PhaseId
	IsSkipped bool
	TargetId  PlayerId
}
