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

type UserRepository struct {
	CollectionName string
	Collection     *mongo.Collection
	Context        context.Context
	RequestID      string
}

func NewUserRepository(requestid string) UserRepository {
	collectionName := "Users"
	c, err := GetMongoDbCollection(collectionName)

	if err != nil {
		log.Fatalln(err)
	}

	return UserRepository{
		CollectionName: collectionName,
		Collection:     c,
		Context:        context.Background(),
		RequestID:      requestid,
	}
}

func (ur UserRepository) FindAll() ([]models.User, error) {
	if ur.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(ur.RequestID, "db:users:find-all")
	defer tracing.End(span)

	results := make([]models.User, 0)
	cur, err := ur.Collection.Find(ur.Context, bson.M{})
	if err != nil {
		return nil, err
	}

	cur.All(ur.Context, &results)

	return results, err
}

func (ur UserRepository) FindAllBy(filter bson.M) ([]models.User, error) {
	span, _ := tracing.Start(ur.RequestID, "db:users:find-all-by", opentracing.Tag{Key: "filter", Value: filter})
	defer tracing.End(span)

	var fetchedUser []models.User

	cur, err := ur.Collection.Find(ur.Context, filter)
	if err != nil {
		return nil, err
	}

	cur.All(ur.Context, &fetchedUser)

	return fetchedUser, nil
}

func (ur UserRepository) FindOneByIdentity(identity string) (*models.User, error) {
	if ur.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(ur.RequestID, "db:users:find-one-by-identity", opentracing.Tag{Key: "identity", Value: identity})
	defer tracing.End(span)

	filter := bson.M{"identity": identity}
	var fetchedUser models.User

	err := ur.Collection.FindOne(ur.Context, filter).Decode(&fetchedUser)
	if err != nil {
		return nil, err
	}

	return &fetchedUser, nil
}

func (ur UserRepository) FindOneById(id string) (*models.User, error) {
	if ur.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(ur.RequestID, "db:users:find-one-by-id", opentracing.Tag{Key: "id", Value: id})
	defer tracing.End(span)

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	var fetchedUser models.User

	err := ur.Collection.FindOne(ur.Context, filter).Decode(&fetchedUser)
	if err != nil {
		return nil, err
	}

	return &fetchedUser, nil
}

func (ur UserRepository) Create(u *models.User) (*models.User, error) {
	if ur.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(ur.RequestID, "db:users:create", opentracing.Tag{Key: "user", Value: *u})
	defer tracing.End(span)

	u.Id = primitive.NewObjectID()

	_, err := ur.Collection.InsertOne(ur.Context, u)

	return u, err
}

func (ur UserRepository) Update(id string, u models.UpdateUserRequest) (*models.User, error) {
	if ur.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(ur.RequestID, "db:users:update",
		opentracing.Tag{Key: "user", Value: id},
		opentracing.Tag{Key: "updated", Value: u},
	)
	defer tracing.End(span)

	user, err := ur.FindOneById(id)
	if err != nil {
		return nil, err
	}

	if u.Name != "" {
		user.Name = u.Name
	}
	if u.Identity != "" {
		user.Identity = u.Identity
	}
	user.Authorization = u.Authorization

	update := bson.M{
		"$set": user,
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err = ur.Collection.UpdateOne(ur.Context, bson.M{"_id": objID}, update)

	return user, err
}

func (ur UserRepository) Delete(id string) error {
	if ur.Collection == nil {
		return errors.New("missing connection")
	}

	span, _ := tracing.Start(ur.RequestID, "db:users:delete", opentracing.Tag{Key: "id", Value: id})
	defer tracing.End(span)

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := ur.Collection.DeleteOne(ur.Context, bson.M{"_id": objID})

	return err
}
