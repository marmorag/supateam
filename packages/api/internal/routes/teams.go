package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/middleware/auth"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
	"log"
)

type TeamRouteHandler struct{}

func (TeamRouteHandler) Register(app fiber.Router) {
	teamsApi := app.Group("/teams")

	teamsApi.Get("/", auth.Authenticated(), getTeams)
	teamsApi.Get("/:id", auth.Authenticated(), getTeam)
	teamsApi.Post("", auth.Authenticated(), createTeam)
	teamsApi.Put("/:id", auth.Authenticated(), updateTeam)
	teamsApi.Delete("/:id", auth.Authenticated(), deleteTeam)

	log.Println("Registered teams api group.")
}

// getTeams godoc
// @Summary List teams
// @Description Get all teams
// @Tags teams
// @Accept  json
// @Produce  json
// @Success 200 array []models.Team
// @Router /teams [get]
func getTeams(c *fiber.Ctx) error {
	tr := repository.NewTeamRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))
	teams, err := tr.FindAll()
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"collection": teams,
		"length":     len(teams),
	})
}

// getTeam godoc
// @Summary Show a team
// @Description get string by ID
// @Tags teams
// @Accept  json
// @Produce  json
// @Param id path string true "Team ID"
// @Success 200 {object} models.Team
// @Router /teams/{id} [get]
func getTeam(c *fiber.Ctx) error {
	tr := repository.NewTeamRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))

	team, err := tr.FindOneById(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(team)
}

// createTeam godoc
// @Summary Create a new team
// @Description Create a new team
// @Tags teams
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Team
// @Router /teams [post]
func createTeam(c *fiber.Ctx) error {
	tr := repository.NewTeamRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))

	createTeamRequest := new(models.CreateTeamRequest)
	if err := c.BodyParser(createTeamRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*createTeamRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	team := models.Team{
		Name:    createTeamRequest.Name,
		Members: createTeamRequest.Members,
	}

	u, err := tr.Create(&team)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(u)
}

// updateTeam godoc
// @Summary Update an existing team
// @Description Update an existing team
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} models.Team
// @Param id path string true "Team ID"
// @Router /teams/{id} [put]
func updateTeam(c *fiber.Ctx) error {
	tr := repository.NewTeamRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))

	updateTeamRequest := new(models.UpdateTeamRequest)
	if err := c.BodyParser(updateTeamRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*updateTeamRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	team, err := tr.Update(c.Params("id"), *updateTeamRequest)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(team)
}

// deleteTeam godoc
// @Summary Delete an existing team
// @Description Delete an existing team
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} models.Team
// @Param id path string true "Team ID"
// @Router /teams/{id} [delete]
func deleteTeam(c *fiber.Ctx) error {
	tr := repository.NewTeamRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))

	err := tr.Delete(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
