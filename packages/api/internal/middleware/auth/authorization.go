package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/tracing"
	"github.com/opentracing/opentracing-go"
)

type ApiGroups string
type ApiAction string

const (
	EventsApiGroup         ApiGroups = "events"
	UsersApiGroup          ApiGroups = "users"
	TeamsApiGroup          ApiGroups = "teams"
	ParticipationsApiGroup ApiGroups = "participations"
)

func (g ApiGroups) S() string {
	return string(g)
}

const (
	AllAction        ApiAction = "*"
	WriteAction      ApiAction = "write"
	UpdateAction     ApiAction = "update"
	DeleteAction     ApiAction = "delete"
	ReadAction       ApiAction = "read"
	WriteSelfAction  ApiAction = "write:self"
	UpdateSelfAction ApiAction = "update:self"
	ReadSelfAction   ApiAction = "read:self"
)

func (a ApiAction) S() string {
	return string(a)
}

func Authorized(api ApiGroups, action ApiAction) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestId := ctx.Locals(internal.GetConfig().RequestIDKey).(string)
		span, _ := tracing.Start(requestId, "auth:authorize",
			opentracing.Tag{Key: string(api), Value: string(action)},
		)

		token := ctx.Locals("user").(*jwt.Token)
		claims := token.Claims.(*ApplicationClaim)
		userAuthorization := claims.UserAuthorization

		span.SetTag("claims", userAuthorization)
		span.SetTag("user", claims.UserId)

		if contains(userAuthorization[api], action) || contains(userAuthorization[api], "*") {
			span.Finish()
			return ctx.Next()
		}

		span.Finish()
		return ctx.SendStatus(fiber.StatusForbidden)
	}
}

func contains(a []ApiAction, x ApiAction) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
