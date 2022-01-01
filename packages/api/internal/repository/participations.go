package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/tracing"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type ParticipationRepository struct {
	CollectionName string
	Collection     *mongo.Collection
	Context        context.Context
	RequestID      string
}

type ResponseFormat string

const (
	ParticipationResponseFormatLong  = "long"
	ParticipationResponseFormatShort = "short"
)

func NewParticipationRepository(requestid string) ParticipationRepository {
	collectionName := "Participations"
	c, err := GetMongoDbCollection(collectionName)

	if err != nil {
		log.Fatalln(err)
	}

	return ParticipationRepository{
		CollectionName: collectionName,
		Collection:     c,
		Context:        context.Background(),
		RequestID:      requestid,
	}
}

func (pr ParticipationRepository) FindAll() ([]models.Participation, error) {
	if pr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(pr.RequestID, "db:participations:create")
	defer tracing.End(span)

	results := make([]models.Participation, 0)
	cur, err := pr.Collection.Find(pr.Context, bson.M{})
	if err != nil {
		return nil, err
	}

	cur.All(pr.Context, &results)

	return results, err
}

// FindAllBy @todo caca, replace with generic when 1.18
func (pr ParticipationRepository) FindAllBy(filter bson.M, format ResponseFormat) ([]models.Participation, []models.ParticipationLong, error) {
	if pr.Collection == nil {
		return nil, nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(pr.RequestID, "db:participations:find-all-by",
		opentracing.Tag{Key: "filter", Value: filter},
		opentracing.Tag{Key: "format", Value: format},
	)
	defer tracing.End(span)

	var fetchedParticipation []models.Participation
	var fetchedParticipationLong []models.ParticipationLong
	aggregatePipeline := make([]bson.M, 0)
	aggregatePipeline = append(aggregatePipeline, bson.M{"$match": filter})

	if format == ParticipationResponseFormatLong {
		aggregatePipeline = append(aggregatePipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         "Users",
					"localField":   "player",
					"foreignField": "_id",
					"as":           "player",
				},
			}, bson.M{
				"$lookup": bson.M{
					"from":         "Teams",
					"localField":   "team",
					"foreignField": "_id",
					"as":           "team",
				},
			})
	}

	cur, err := pr.Collection.Aggregate(pr.Context, aggregatePipeline)
	if err != nil {
		return nil, nil, err
	}

	if format == ParticipationResponseFormatLong {
		cur.All(pr.Context, &fetchedParticipationLong)
		return nil, fetchedParticipationLong, nil
	}

	cur.All(pr.Context, &fetchedParticipation)
	return fetchedParticipation, nil, nil
}

func (pr ParticipationRepository) FindOneById(id string) (*models.Participation, error) {
	if pr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(pr.RequestID, "db:participations:find-one-by-id", opentracing.Tag{Key: "id", Value: id})
	defer tracing.End(span)

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	var fetchedParticipation models.Participation

	err := pr.Collection.FindOne(pr.Context, filter).Decode(&fetchedParticipation)
	if err != nil {
		return nil, err
	}

	return &fetchedParticipation, nil
}

func (pr ParticipationRepository) Create(p *models.Participation) (*models.Participation, error) {
	if pr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(pr.RequestID, "db:participations:create", opentracing.Tag{Key: "participation", Value: *p})
	defer tracing.End(span)

	p.Id = primitive.NewObjectID()

	_, err := pr.Collection.InsertOne(pr.Context, p)

	return p, err
}

func (pr ParticipationRepository) Update(id string, p models.UpdateParticipationRequest) (*models.Participation, error) {
	if pr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(pr.RequestID, "db:participations:update",
		opentracing.Tag{Key: "id", Value: id},
		opentracing.Tag{Key: "participation", Value: p},
	)
	defer tracing.End(span)

	participation, err := pr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	participation.Status = p.Status

	update := bson.M{
		"$set": participation,
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err = pr.Collection.UpdateOne(pr.Context, bson.M{"_id": objID}, update)

	return participation, err
}

func (pr ParticipationRepository) Delete(id string) error {
	if pr.Collection == nil {
		return errors.New("missing connection")
	}

	span, _ := tracing.Start(pr.RequestID, "db:participations:delete", opentracing.Tag{Key: "id", Value: id})
	defer tracing.End(span)

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := pr.Collection.DeleteOne(pr.Context, bson.M{"_id": objID})

	return err
}

func (pr ParticipationRepository) deleteBatch(ps []models.Participation) error {
	span, _ := tracing.Start(pr.RequestID, "db:participations:delete-batch", opentracing.Tag{Key: "participations", Value: len(ps)})
	defer tracing.End(span)

	pIDs := make([]primitive.ObjectID, 0)
	for _, p := range ps {
		pIDs = append(pIDs, p.Id)
	}

	_, err := pr.Collection.DeleteMany(pr.Context, bson.M{
		"_id": bson.M{
			"$in": pIDs,
		},
	})

	return err
}

func (pr ParticipationRepository) createBatch(ps []models.Participation) error {
	span, _ := tracing.Start(pr.RequestID, "db:participations:create-batch", opentracing.Tag{Key: "participations", Value: len(ps)})
	defer tracing.End(span)

	iPs := make([]interface{}, 0)
	for _, p := range ps {
		p.Id = primitive.NewObjectID()
		iPs = append(iPs, p)
	}

	_, err := pr.Collection.InsertMany(pr.Context, iPs)
	return err
}

func (pr ParticipationRepository) SyncParticipation(event *models.Event) error {
	if pr.Collection == nil {
		return errors.New("missing connection")
	}

	span, _ := tracing.Start(pr.RequestID, "db:participations:synchronize", opentracing.Tag{Key: "event", Value: *event})
	defer tracing.End(span)

	tr := NewTeamRepository(pr.RequestID)

	participations, _, err := pr.FindAllBy(bson.M{
		"event": event.Id,
	}, ParticipationResponseFormatShort)

	if err != nil {
		return errors.New(fmt.Sprintf("unable to fetch participation : %s", err.Error()))
	}

	partToRemove := make([]models.Participation, 0)
	for _, participation := range participations {
		// handle to remove for Team & Players
		if !event.HasParticipation(participation) {
			partToRemove = append(partToRemove, participation)
		}
	}

	partToAdd := make([]models.Participation, 0)
	for _, teamID := range event.Teams {
		if !models.IncludeObject(participations, teamID) {
			team, err := tr.FindOneById(teamID.Hex())
			if err != nil {
				return errors.New(fmt.Sprintf("unable to fetch team : %s", err.Error()))
			}

			for _, player := range team.Members {
				partToAdd = append(partToAdd, models.Participation{
					Event:  event.Id,
					Player: player,
					Team:   teamID,
					Status: models.ParticipationUnknown,
				})
			}
		}
	}

	for _, playerID := range event.Players {
		// only add player if it does not already exist through a team
		if !models.IncludeObject(participations, playerID) && !models.IncludeObject(partToAdd, playerID) {
			partToAdd = append(partToAdd, models.Participation{
				Event:  event.Id,
				Player: playerID,
				Team:   primitive.NilObjectID,
				Status: models.ParticipationUnknown,
			})
		}
	}

	if len(partToRemove) > 0 {
		err = pr.deleteBatch(partToRemove)
		if err != nil {
			return errors.New(fmt.Sprintf("unable to delete participations : %s", err.Error()))
		}
	}

	if len(partToAdd) > 0 {
		err = pr.createBatch(partToAdd)
		if err != nil {
			return errors.New(fmt.Sprintf("unable to create participations : %s", err.Error()))
		}
	}

	return nil
}
