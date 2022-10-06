package middleware

import (
	"fmt"
	"uwwolf/app/enum"
	"uwwolf/app/validator"

	"github.com/gofiber/fiber/v2"
)

func Validation[T any](c *fiber.Ctx) error {
	parsedData := new(T)

	if err := c.BodyParser(parsedData); err != nil {
		fmt.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": enum.InvalidInputErrorTag,
			"error":   "Unknown error!",
		})
	}

	if err := validator.ValidateStruct(parsedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": enum.InvalidInputErrorTag,
			"error":   err,
		})
	}

	c.Locals(enum.FiberLocalDataKey, parsedData)

	return c.Next()
}
