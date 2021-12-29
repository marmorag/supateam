package repository

import (
	"context"
	"errors"
	"github.com/marmorag/supateam/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type TeamRepository struct {
	CollectionName string
	Collection     *mongo.Collection
	Context        context.Context
}

func NewTeamRepository() TeamRepository {
	collectionName := "Teams"
	c, err := GetMongoDbCollection(collectionName)

	if err != nil {
		log.Fatalln(err)
	}

	return TeamRepository{
		CollectionName: collectionName,
		Collection:     c,
		Context:        context.Background(),
	}
}

func (tr TeamRepository) FindAll() ([]models.Team, error) {
	if tr.Collection == nil {
		return nil, errors.New("missing connection")
	}

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

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	var fetchedTeam models.Team

	err := tr.Collection.FindOne(tr.Context, filter).Decode(&fetchedTeam)
	if err != nil {
		return nil, err
	}

	return &fetchedTeam, nil
}

func (tr TeamRepository) Create(e *models.Team) (*models.Team, error) {
	if tr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	e.Id = primitive.NewObjectID()

	_, err := tr.Collection.InsertOne(tr.Context, e)

	return e, err
}

func (tr TeamRepository) Update(id string, e models.UpdateTeamRequest) (*models.Team, error) {
	if tr.Collection == nil {
		return nil, errors.New("missing connection")
	}

	team, err := tr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	team.Name = e.Name
	team.Members = e.Members

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

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := tr.Collection.DeleteOne(tr.Context, bson.M{"_id": objID})

	return err
}
