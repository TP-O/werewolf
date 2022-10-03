package enum

import (
	"uwwolf/app/types"
)

const (
	UnlimitedTimes types.Expiration = iota - 1
	OutOfTimes
	OneTimes
)
