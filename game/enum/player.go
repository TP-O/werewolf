package enum

type PlayerStatus = uint8

const (
	OfflineStatus PlayerStatus = iota
	OnlineStatus
	BusyStatus
	InGameStatus
)

type PlayerID = string

func IsUnknownPlayerID(p PlayerID) bool {
	return p == ""
}
