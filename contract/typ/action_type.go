package typ

type ActionInstruction struct {
	GameId  string `validate:"required,len=20,alphanum"`
	Actor   uint   `validate:"required,min=1"`
	Targets []uint `validate:"required_if=Skipped false,min=1,dive,min=1"`
	Skipped bool   `validate:""`
	Payload []byte `validate:"required"`
}
