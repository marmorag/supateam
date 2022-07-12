package testing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/marmorag/supateam/internal"
)

func BuildTestApplication() (*fiber.App, *fiber.Router) {
	config := internal.GetConfig()

	app := fiber.New(fiber.Config{
		Prefork: config.ApplicationPrefork,
		AppName: config.ApplicationName,
	})

	// RequestID Middleware - Add request id header & local to ctx + response
	app.Use(requestid.New())

	api := app.Group("/api")

	return app, &api
}
