package server

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"uwwolf/app/handler"
	"uwwolf/app/middleware"
	"uwwolf/app/types"
	"uwwolf/config"
)

func StartAPI() {
	app := fiber.New()

	app.Post("/api/v1/start", middleware.Validation[types.GameSetting], handler.StartGame)

	log.Fatal(app.Listen(":" + strconv.Itoa(config.App.HttpPort)))
}
