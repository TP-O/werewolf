package types

type ActionData struct {
	GameId  string `validate:"required,len=20,alphanum"`
	Actor   int    `validate:"required,min=1"`
	Targets []int  `validate:"required_if=Skipped false,min=1,dive,min=1"`
	Skipped bool   `validate:""`
	Payload []byte `validate:"required"`
}
