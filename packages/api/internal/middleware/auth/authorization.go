package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/tracing"
	"github.com/opentracing/opentracing-go"
	"regexp"
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

func isSelfApiAction(action ApiAction) bool {
	selfActions := []ApiAction{
		WriteSelfAction,
		UpdateSelfAction,
		ReadSelfAction,
	}

	return contains(selfActions, action)
}

func getElevated(action ApiAction) ApiAction {
	switch action {
	case WriteSelfAction:
		return WriteAction
	case UpdateSelfAction:
		return UpdateAction
	case ReadSelfAction:
		return ReadAction
	default:
		return AllAction
	}
}

func (a ApiAction) S() string {
	return string(a)
}

type SelfActionHandler interface {
	Vote(userId string, entityId string) bool
}

func Authorized(api ApiGroups, action ApiAction, handlers ...SelfActionHandler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestId := ctx.Locals(internal.GetConfig().RequestIDKey).(string)
		span, _ := tracing.Start(requestId, "auth:authorize",
			opentracing.Tag{Key: string(api), Value: string(action)},
		)

		token := ctx.Locals("user").(*jwt.Token)
		claims := token.Claims.(*ApplicationClaim)

		tracing.SetTag(span, "claims", claims)
		tracing.SetTag(span, "user", claims.UserId)

		if enforced, _ := enforce(*claims, api, action, ctx, handlers); enforced {
			tracing.End(span)
			return ctx.Next()
		}

		tracing.End(span)
		return ctx.SendStatus(fiber.StatusForbidden)
	}
}

func enforce(claims ApplicationClaim, api ApiGroups, action ApiAction, ctx *fiber.Ctx, handlers []SelfActionHandler) (bool, error) {
	if elevated := getElevated(action); contains(claims.UserAuthorization[api], elevated) || contains(claims.UserAuthorization[api], AllAction) {
		return true, nil
	}

	if isSelfApiAction(action) {

		if len(handlers) == 0 {
			return false, errors.New("no handlers for a self managed entity")
		}

		entityId := extractEntityId(ctx)

		isEnforced := false
		for _, handler := range handlers {
			isEnforced = isEnforced || handler.Vote(claims.UserId, entityId)
		}
		return isEnforced, nil
	}

	return contains(claims.UserAuthorization[api], action), nil
}

func extractEntityId(ctx *fiber.Ctx) string {
	rg := regexp.MustCompile(`([[:xdigit:]]{24})`)
	return rg.FindString(ctx.Path())
}

func contains(a []ApiAction, x ApiAction) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
