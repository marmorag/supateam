package repository

import (
	"context"
	"errors"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/tracing"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type TeamRepository struct {
	CollectionName string
	Collection     *mongo.Collection
	Context        context.Context
	RequestID      string
}

func NewTeamRepository(requestid string) TeamRepository {
	collectionName := "Teams"
	c, err := GetMongoDbCollection(collectionName)

	if err != nil {
		log.Fatalln(err)
	}

	return TeamRepository{
		CollectionName: collectionName,
		Collection:     c,
		Context:        context.Background(),
		RequestID:      requestid,
	}
}

func (tr TeamRepository) FindAll() ([]models.Team, error) {
	if tr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(tr.RequestID, "db:teams:find-all")
	defer tracing.End(span)

	results := make([]models.Team, 0)
	cur, err := tr.Collection.Find(tr.Context, bson.M{})
	if err != nil {
		return nil, err
	}

	cur.All(tr.Context, &results)

	return results, err
}

func (tr TeamRepository) FindAllBy(filter bson.M) ([]models.Team, error) {
	if tr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(tr.RequestID, "db:teams:find-all-by", opentracing.Tag{Key: "filter", Value: filter})
	defer tracing.End(span)

	var fetchedTeam []models.Team

	cur, err := tr.Collection.Find(tr.Context, filter)
	if err != nil {
		return nil, err
	}

	cur.All(tr.Context, &fetchedTeam)

	return fetchedTeam, nil
}

func (tr TeamRepository) FindOneById(id string) (*models.Team, error) {
	if tr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(tr.RequestID, "db:teams:find-one-by-id", opentracing.Tag{Key: "id", Value: id})
	defer tracing.End(span)

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	var fetchedTeam models.Team

	err := tr.Collection.FindOne(tr.Context, filter).Decode(&fetchedTeam)
	if err != nil {
		return nil, err
	}

	return &fetchedTeam, nil
}

func (tr TeamRepository) Create(t *models.Team) (*models.Team, error) {
	if tr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(tr.RequestID, "db:teams:create", opentracing.Tag{Key: "team", Value: *t})
	defer tracing.End(span)

	t.Id = primitive.NewObjectID()

	_, err := tr.Collection.InsertOne(tr.Context, t)

	return t, err
}

func (tr TeamRepository) Update(id string, t models.UpdateTeamRequest) (*models.Team, error) {
	if tr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(tr.RequestID, "db:teams:update",
		opentracing.Tag{Key: "id", Value: id},
		opentracing.Tag{Key: "team", Value: t},
	)
	defer tracing.End(span)

	team, err := tr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	team.Name = t.Name
	team.Members = t.Members

	update := bson.M{
		"$set": team,
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err = tr.Collection.UpdateOne(tr.Context, bson.M{"_id": objID}, update)

	return team, err
}

func (tr TeamRepository) Delete(id string) error {
	if tr.Collection == nil {
		return errors.New("missing connection")
	}

	span, _ := tracing.Start(tr.RequestID, "db:teams:delete", opentracing.Tag{Key: "id", Value: id})
	defer tracing.End(span)

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := tr.Collection.DeleteOne(tr.Context, bson.M{"_id": objID})

	return err
}
