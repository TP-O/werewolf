package types

type PollRound = map[PlayerId]*PollRecord

type PollRecord struct {
	ElectorIds []PlayerId `json:"electorIds"`
	Votes      uint       `json:"votes"`
}
