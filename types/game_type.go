package types

type GameInstance struct {
	GameId             string `validate:"required,len=20,alphanum"`
	Capacity           int    `validate:"required,game_capacity"`
	NumberOfWerewolves int    `validate:"required,min=1"`
	RolePool           []int  `validate:"required,unique,dive,min=3"`
}
