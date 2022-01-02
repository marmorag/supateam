package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/middleware/auth"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
	"github.com/marmorag/supateam/internal/tracing"
	"log"
)

type ParticipationRouteHandler struct{}

func (ParticipationRouteHandler) Register(app fiber.Router) {
	participationsApi := app.Group("/participations")

	participationsApi.Get("/",
		auth.Authenticated(),
		tracing.HandlerTracer("get-participations"),
		getParticipations,
	)
	participationsApi.Get("/:id",
		auth.Authenticated(),
		tracing.HandlerTracer("get-participation"),
		getParticipation,
	)
	participationsApi.Post("",
		auth.Authenticated(),
		auth.Authorized(auth.ParticipationsApiGroup, auth.WriteSelfAction),
		tracing.HandlerTracer("create-participation"),
		createParticipation,
	)
	participationsApi.Put("/:id",
		auth.Authenticated(),
		auth.Authorized(auth.ParticipationsApiGroup, auth.UpdateSelfAction),
		tracing.HandlerTracer("update-participation"),
		updateParticipation,
	)
	participationsApi.Delete("/:id",
		auth.Authenticated(),
		auth.Authorized(auth.ParticipationsApiGroup, auth.DeleteAction),
		tracing.HandlerTracer("delete-participation"),
		deleteParticipation,
	)

	log.Println("Registered participations api group.")
}

func (h ParticipationRouteHandler) Vote(userId string, entityId string) bool {
	return true
}

// getParticipations godoc
// @Summary List participations
// @Description Get all participations
// @Tags participations
// @Accept  json
// @Produce  json
// @Success 200 array []models.Participation
// @Router /participations [get]
func getParticipations(c *fiber.Ctx) error {
	pr := repository.NewParticipationRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))
	participations, err := pr.FindAll()
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"collection": participations,
		"length":     len(participations),
	})
}

// getParticipation godoc
// @Summary Show a participation
// @Description get string by ID
// @Tags participations
// @Accept  json
// @Produce  json
// @Param id path string true "Participation ID"
// @Success 200 {object} models.Participation
// @Router /participations/{id} [get]
func getParticipation(c *fiber.Ctx) error {
	pr := repository.NewParticipationRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))

	participation, err := pr.FindOneById(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(participation)
}

// createParticipation godoc
// @Summary Create a new participation
// @Description Create a new participation
// @Tags participations
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Participation
// @Router /participations [post]
func createParticipation(c *fiber.Ctx) error {
	pr := repository.NewParticipationRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))

	createParticipationRequest := new(models.CreateParticipationRequest)
	if err := c.BodyParser(createParticipationRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*createParticipationRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	participation := models.Participation{
		Status: createParticipationRequest.Status,
		Event:  createParticipationRequest.Event,
		Player: createParticipationRequest.Player,
		Team:   createParticipationRequest.Team,
	}

	u, err := pr.Create(&participation)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(u)
}

// updateParticipation godoc
// @Summary Update an existing participation
// @Description Update an existing participation
// @Tags participations
// @Accept json
// @Produce json
// @Success 200 {object} models.Participation
// @Param id path string true "Participation ID"
// @Router /participations/{id} [put]
func updateParticipation(c *fiber.Ctx) error {
	pr := repository.NewParticipationRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))

	updateParticipationRequest := new(models.UpdateParticipationRequest)
	if err := c.BodyParser(updateParticipationRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*updateParticipationRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	participation, err := pr.Update(c.Params("id"), *updateParticipationRequest)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(participation)
}

// deleteParticipation godoc
// @Summary Delete an existing participation
// @Description Delete an existing participation
// @Tags participations
// @Accept json
// @Produce json
// @Success 200 {object} models.Participation
// @Param id path string true "Participation ID"
// @Router /participations/{id} [delete]
func deleteParticipation(c *fiber.Ctx) error {
	pr := repository.NewParticipationRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))

	err := pr.Delete(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
