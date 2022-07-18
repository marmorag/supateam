package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
)

type ParticipationPostActionHandler struct{}

func (p ParticipationPostActionHandler) Vote(ctx *fiber.Ctx, userId string, entityId string) bool {
	createParticipationRequest := new(models.CreateParticipationRequest)
	if err := ctx.BodyParser(createParticipationRequest); err != nil {
		return false
	}

	return createParticipationRequest.Player.Hex() == userId
}

func (p ParticipationPostActionHandler) ExtractEntityId(ctx *fiber.Ctx) (string, error) {
	// Return empty entity ID as we will handle it later in Vote method
	return "", nil
}

type ParticipationPutActionHandler struct{}

func (p ParticipationPutActionHandler) Vote(ctx *fiber.Ctx, userId string, entityId string) bool {
	requestId := ctx.Locals(internal.GetConfig().RequestIDKey).(string)
	er := repository.NewParticipationRepository(requestId)

	participation, err := er.FindOneById(entityId)
	if err != nil {
		log.Printf("found empty record for entity %s", entityId)
		return false
	}

	uId, _ := primitive.ObjectIDFromHex(userId)
	log.Printf("found participation %s with authenticated %s", participation.Player, uId)
	return participation.Player == uId
}

func (p ParticipationPutActionHandler) ExtractEntityId(ctx *fiber.Ctx) (string, error) {
	rg := regexp.MustCompile(`([[:xdigit:]]{24})`)
	return rg.FindString(ctx.Path()), nil
}
