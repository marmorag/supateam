package routes

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
)

type RouteHandler interface {
	Register(app fiber.Router)
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateRequest(request interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func GetUserFromContext(c *fiber.Ctx) (*models.User, error) {
	token := c.Locals("user").(*jwt.Token)

	userId := models.GetUserIdFromToken(token)
	ur := repository.NewUserRepository()

	return ur.FindOneById(userId)
}

func jsonError(c *fiber.Ctx, errorCode int, errorMessage interface{}) error {
	return c.Status(errorCode).JSON(fiber.Map{
		"message": errorMessage,
	})
}
