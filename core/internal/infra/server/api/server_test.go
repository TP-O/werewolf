package api_test

// import (
// 	"testing"
// 	"uwwolf/internal/app/game/logic/types"

// 	"github.com/gin-gonic/gin/binding"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/stretchr/testify/suite"
// )

// type ApiServiceSuite struct {
// 	suite.Suite
// 	playerID1 types.PlayerID
// 	playerID2 types.PlayerID
// }

// func (ass *ApiServiceSuite) SetupSuite() {
// 	ass.playerID1 = types.PlayerID("1")
// 	ass.playerID2 = types.PlayerID("2")

// 	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
// 		validation.Setup(v)
// 	}
// }

// func TestGameServiceSuite(t *testing.T) {
// 	suite.Run(t, new(ApiServiceSuite))
// }
