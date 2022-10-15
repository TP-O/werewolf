package types

type VoteRecord struct {
	ElectorIds []PlayerId `json:"electorIds"`
	Votes      uint       `json:"votes"`
}

type PollRecord struct {
	Votes    map[PlayerId]*VoteRecord `json:"votes"`
	Winner   PlayerId                 `json:"winner"`
	IsClosed bool                     `json:"isClosed"`
}
