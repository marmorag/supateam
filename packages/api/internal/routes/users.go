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

type UserRouteHandler struct{}

func (UserRouteHandler) Register(app fiber.Router) {
	usersApi := app.Group("/users")

	usersApi.Get("/", getUsers)
	usersApi.Get("/:id", getUser)
	usersApi.Get("/:id/participations", getUserParticipation)
	usersApi.Post("", createUser)
	usersApi.Put("/:id", auth.Authenticated(), updateUser)
	usersApi.Delete("/:id", auth.Authenticated(), deleteUser)

	log.Println("Registered users api group.")
}

// getUsers godoc
// @Summary List users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 array []models.User
// @Router /users [get]
func getUsers(c *fiber.Ctx) error {
	ur := repository.NewUserRepository()
	users, err := ur.FindAll()
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"collection": users,
		"length":     len(users),
	})
}

// getUser godoc
// @Summary Show a user
// @Description get string by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func getUser(c *fiber.Ctx) error {
	ur := repository.NewUserRepository()

	user, err := ur.FindOneById(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(user)
}

// getUserParticipation godoc
// @Summary Show User participations
// @Description get string by ID
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} []models.Participation
// @Router /users/{id}/participations [get]
func getUserParticipation(c *fiber.Ctx) error {
	er := repository.NewParticipationRepository()

	responseFormat := repository.ResponseFormat(c.Query("format", repository.ParticipationResponseFormatShort))

	userOID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	participations, longParticipations, err := er.FindAllBy(bson.M{
		"player": userOID,
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

// createUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Router /users [post]
func createUser(c *fiber.Ctx) error {
	ur := repository.NewUserRepository()

	createUserRequest := new(models.CreateUserRequest)
	if err := c.BodyParser(createUserRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*createUserRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	user := models.User{
		Name:          createUserRequest.Name,
		Identity:      createUserRequest.Identity,
		Authorization: createUserRequest.Authorization,
	}

	u, err := ur.Create(&user)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(u)
}

// updateUser godoc
// @Summary Update an existing user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Param id path string true "User ID"
// @Router /users/{id} [put]
func updateUser(c *fiber.Ctx) error {
	ur := repository.NewUserRepository()

	updateUserRequest := new(models.UpdateUserRequest)
	if err := c.BodyParser(updateUserRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*updateUserRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	user, err := ur.Update(c.Params("id"), *updateUserRequest)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(user)
}

// deleteUser godoc
// @Summary Delete an existing user
// @Description Delete an existing user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Param id path string true "User ID"
// @Router /users/{id} [delete]
func deleteUser(c *fiber.Ctx) error {
	ur := repository.NewUserRepository()

	err := ur.Delete(c.Params("id"))
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
