package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal/middleware/auth"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type EventRouteHandler struct{}

func (EventRouteHandler) Register(app fiber.Router) {
	eventsApi := app.Group("/events")

	eventsApi.Get("/", getEvents)
	eventsApi.Get("/:id", getEvent)
	eventsApi.Get("/:id/participations", getEventParticipation)
	eventsApi.Post("", createEvent)
	eventsApi.Put("/:id", auth.Authenticated(), updateEvent)
	eventsApi.Delete("/:id", auth.Authenticated(), deleteEvent)

	log.Println("Registered events api group.")
}

// getEvents godoc
// @Summary List events
// @Description Get all events
// @Tags events
// @Accept  json
// @Produce  json
// @Success 200 array []models.Event
// @Router /events [get]
func getEvents(c *fiber.Ctx) error {
	er := repository.NewEventRepository()
	events, err := er.FindAll()
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"collection": events,
		"length":     len(events),
	})
}

// getEvent godoc
// @Summary Show as event
// @Description get string by ID
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path string true "Event ID"
// @Success 200 {object} models.Event
// @Router /events/{id} [get]
func getEvent(c *fiber.Ctx) error {
	er := repository.NewEventRepository()

	event, err := er.FindOneById(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(event)
}

// getEventParticipation godoc
// @Summary Show event participations
// @Description get string by ID
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path string true "Event ID"
// @Success 200 {object} []models.Participation
// @Router /events/{id}/participations [get]
func getEventParticipation(c *fiber.Ctx) error {
	er := repository.NewParticipationRepository()

	responseFormat := repository.ResponseFormat(c.Query("format", repository.ParticipationResponseFormatShort))

	eventOID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	participations, longParticipations, err := er.FindAllBy(bson.M{
		"event": eventOID,
	}, responseFormat)

	if err != nil {
		return jsonError(c, fiber.StatusNotFound, err.Error())
	}

	if responseFormat == repository.ParticipationResponseFormatLong {
		return c.JSON(fiber.Map{
			"collection": longParticipations,
			"length":     len(participations),
		})
	}

	return c.JSON(fiber.Map{
		"collection": participations,
		"length":     len(participations),
	})
}

// createEvent godoc
// @Summary Create a new event
// @Description Create a new event
// @Tags events
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Event
// @Router /events [post]
func createEvent(c *fiber.Ctx) error {
	er := repository.NewEventRepository()
	pr := repository.NewParticipationRepository()

	createEventRequest := new(models.CreateEventRequest)
	if err := c.BodyParser(createEventRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*createEventRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	event := models.Event{
		Title:       createEventRequest.Title,
		Date:        createEventRequest.Date,
		Description: createEventRequest.Description,
		Duration:    createEventRequest.Duration,
		Kind:        createEventRequest.Kind,
		Teams:       createEventRequest.Teams,
		Players:     createEventRequest.Players,
	}

	created, err := er.Create(&event)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	err = pr.SyncParticipation(created)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// updateEvent godoc
// @Summary Update an existing event
// @Description Update an existing event
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {object} models.Event
// @Param id path string true "Event ID"
// @Router /events/{id} [put]
func updateEvent(c *fiber.Ctx) error {
	er := repository.NewEventRepository()
	pr := repository.NewParticipationRepository()

	updateEventRequest := new(models.UpdateEventRequest)
	if err := c.BodyParser(updateEventRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*updateEventRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	updated, err := er.Update(c.Params("id"), *updateEventRequest)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	err = pr.SyncParticipation(updated)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(updated)
}

// deleteEvent godoc
// @Summary Delete an existing event
// @Description Delete an existing event
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {object} models.Event
// @Param id path string true "Event ID"
// @Router /events/{id} [delete]
func deleteEvent(c *fiber.Ctx) error {
	er := repository.NewEventRepository()

	err := er.Delete(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}