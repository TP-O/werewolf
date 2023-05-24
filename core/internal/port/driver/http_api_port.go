package driver

import "github.com/gin-gonic/gin"

type HttpApiPort interface {
	UpdateGameSetting(ctx *gin.Context)
	StartGame(ctx *gin.Context)
}
