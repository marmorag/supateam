package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal"
)

// Authenticated protect routes
func Authenticated() func(*fiber.Ctx) error {
	return New(Config{
		SigningKey:   []byte(internal.GetConfig().ApplicationSecret),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}
