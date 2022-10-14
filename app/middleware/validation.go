package middleware

import (
	"github.com/gofiber/fiber/v2"

	"uwwolf/app/enum"
	"uwwolf/app/validator"
)

func Validation[T any](c *fiber.Ctx) error {
	parsedPayload := new(T)

	if err := c.BodyParser(parsedPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
			"error": fiber.Map{
				"tag":     enum.InvalidInputErrorTag,
				"message": "Invalid payload!",
			},
		})
	}

	if err := validator.ValidateStruct(parsedPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
			"error": fiber.Map{
				"tag":     enum.InvalidInputErrorTag,
				"message": err,
			},
		})
	}

	// Store the parsed payload for use in next handlers
	c.Locals(enum.FiberLocalPayloadKey, parsedPayload)

	return c.Next()
}
