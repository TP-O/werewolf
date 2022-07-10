package contract

type GameInstanceInit struct {
	GameId             string `validate:"required,len=20,alphanum"`
	Capacity           uint   `validate:"required,game_capacity"`
	NumberOfWerewolves uint   `validate:"required,min=1"`
	RolePool           []uint `validate:"required,unique,dive,min=3"`
}
