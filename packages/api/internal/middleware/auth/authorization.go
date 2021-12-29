package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type ApiGroups string
type ApiAction string

const (
	EventsApiGroup ApiGroups = "events"
	UsersApiGroup  ApiGroups = "users"
)

const (
	AllAction   ApiAction = "*"
	WriteAction ApiAction = "write"
	ReadAction  ApiAction = "read"
)

func Authorized(api ApiGroups, action ApiAction) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Locals("user").(*jwt.Token)
		claims := token.Claims.(*ApplicationClaim)
		userAuthorization := claims.UserAuthorization

		if contains(userAuthorization[string(api)], string(action)) || contains(userAuthorization[string(api)], "*") {
			return ctx.Next()
		}

		return ctx.SendStatus(fiber.StatusForbidden)
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
