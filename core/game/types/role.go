package types

// FactionID is ID of faction.
type FactionID uint8

// IsUnknown checks if faction ID is 0.
func (f FactionID) IsUnknown() bool {
	return f == 0
}

// RoleID is ID type of role.
type RoleID uint8

// IsUnknown checks if role ID is 0.
func (r RoleID) IsUnknown() bool {
	return r == 0
}

// ActivateAbilityRequest contains information for ability activating.
type ActivateAbilityRequest struct {
	// AbilityIndex is activated ability index.
	AbilityIndex uint8

	// TargetID  is player ID of target player.
	TargetID PlayerID

	// IsSkipped marks that the request is ignored.
	IsSkipped bool
}
