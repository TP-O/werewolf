package types

type Status uint

type PlayerID string

func (p PlayerID) IsUnknown() bool {
	return p == ""
}
