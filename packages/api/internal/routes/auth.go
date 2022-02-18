package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
	"github.com/marmorag/supateam/internal/tracing"
	"log"
)

type AuthRouteHandler struct{}

func (AuthRouteHandler) Register(app fiber.Router) {
	authApi := app.Group("/auth")

	authApi.Post("/login",
		tracing.HandlerTracer("auth-user"),
		authUser,
	)

	log.Println("Registered auth api group.")
}

type AuthRequest struct {
	Identity string `validate:"required"`
}

// authUser godoc
// @Summary Authenticate a user
// @Description Authenticate a user, receiving identify and password, returning JWT Token
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Router /auth/login [post]
func authUser(c *fiber.Ctx) error {
	authRequest := new(AuthRequest)
	if err := c.BodyParser(authRequest); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	validationErrors := ValidateRequest(*authRequest)
	if validationErrors != nil {
		return jsonError(c, fiber.StatusBadRequest, validationErrors)
	}

	userRepository := repository.NewUserRepository(c.Locals(internal.GetConfig().RequestIDKey).(string))
	user, err := userRepository.FindOneByIdentity(authRequest.Identity)
	if err != nil {
		return jsonError(c, fiber.StatusNotFound, err)
	}

	// Disable password authentication, for simplicity reason only phone number is used to authenticate
	//if !models.CheckPasswordHash(authRequest.Id, user.Identity) {
	//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//		"success": false,
	//		"message": "invalid password",
	//	})
	//}

	token := models.BuildToken(*user)

	t, err := token.SignedString([]byte(internal.GetConfig().ApplicationSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   t,
	})
}
