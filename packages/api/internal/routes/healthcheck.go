package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal/repository"
	"log"
)

type HealthcheckRouteHandler struct{}

func (HealthcheckRouteHandler) Register(app fiber.Router) {
	app.Get("/healthz/alive", handleAlive)
	app.Get("/healthz/ready", handleReady)

	log.Println("Registered healthcheck api group.")
}

// handleSubscribe godoc
// @Summary Healthcheck route to ping API
// @Description Healthcheck route to ping API, will respond OK
// @Tags healthcheck
// @Accept  json
// @Produce  json
// @Success 200
// @Router /healthz/ready [get]
func handleAlive(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

// handleSubscribe godoc
// @Summary Healthcheck route to ping API
// @Description Healthcheck route to ping API, will respond OK if it can connect to mongo
// @Tags healthcheck
// @Accept  json
// @Produce  json
// @Success 200
// @Router /healthz/alive [get]
func handleReady(c *fiber.Ctx) error {
	_ = repository.GetMongoConnection()
	return c.SendStatus(fiber.StatusOK)
}
