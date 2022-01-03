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

type EventRepository struct {
	CollectionName string
	Collection     *mongo.Collection
	Context        context.Context
	RequestID      string
}

func NewEventRepository(requestid string) EventRepository {
	collectionName := "Events"
	c, err := GetMongoDbCollection(collectionName)

	if err != nil {
		log.Fatalln(err)
	}

	return EventRepository{
		CollectionName: collectionName,
		Collection:     c,
		Context:        context.Background(),
		RequestID:      requestid,
	}
}

func (er EventRepository) FindAll() ([]models.Event, error) {
	if er.Collection == nil {
		return nil, errors.New("missing connection")
	}
	span, _ := tracing.Start(er.RequestID, "db:events:find-all")
	defer tracing.End(span)

	results := make([]models.Event, 0)
	cur, err := er.Collection.Find(er.Context, bson.M{})
	if err != nil {
		return nil, err
	}

	cur.All(er.Context, &results)

	return results, err
}

func (er EventRepository) FindAllBy(filter bson.M) ([]models.Event, error) {
	span, _ := tracing.Start(er.RequestID, "db:events:find-all-by", opentracing.Tag{Key: "filter", Value: filter})
	defer tracing.End(span)

	var fetchedEvent []models.Event

	cur, err := er.Collection.Find(er.Context, filter)
	if err != nil {
		return nil, err
	}

	cur.All(er.Context, &fetchedEvent)

	return fetchedEvent, nil
}

func (er EventRepository) FindOneById(id string) (*models.Event, error) {
	if er.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(er.RequestID, "db:events:find-one-by-id", opentracing.Tag{Key: "id", Value: id})
	defer tracing.End(span)

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	var fetchedEvent models.Event

	err := er.Collection.FindOne(er.Context, filter).Decode(&fetchedEvent)
	if err != nil {
		return nil, err
	}

	return &fetchedEvent, nil
}

func (er EventRepository) Create(e *models.Event) (*models.Event, error) {
	if er.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(er.RequestID, "db:events:create", opentracing.Tag{Key: "event", Value: *e})
	defer tracing.End(span)

	e.Id = primitive.NewObjectID()
	arrayGuard(e)

	_, err := er.Collection.InsertOne(er.Context, e)

	return e, err
}

func (er EventRepository) Update(id string, e models.UpdateEventRequest) (*models.Event, error) {
	if er.Collection == nil {
		return nil, errors.New("missing connection")
	}

	span, _ := tracing.Start(er.RequestID, "db:events:update",
		opentracing.Tag{Key: "event", Value: id},
		opentracing.Tag{Key: "updated", Value: e},
	)
	defer tracing.End(span)

	event, err := er.FindOneById(id)
	if err != nil {
		return nil, err
	}

	event.Title = e.Title
	event.Description = e.Description
	event.Date = e.Date
	event.Duration = e.Duration
	event.Kind = e.Kind
	event.Teams = e.Teams
	event.Players = e.Players
	arrayGuard(event)

	update := bson.M{
		"$set": event,
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err = er.Collection.UpdateOne(er.Context, bson.M{"_id": objID}, update)

	return event, err
}

func (er EventRepository) Delete(id string) error {
	if er.Collection == nil {
		return errors.New("missing connection")
	}

	span, _ := tracing.Start(er.RequestID, "db:events:delete", opentracing.Tag{Key: "id", Value: id})
	defer tracing.End(span)

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := er.Collection.DeleteOne(er.Context, bson.M{"_id": objID})

	return err
}

func arrayGuard(event *models.Event) {
	if event.Players == nil {
		event.Players = make([]primitive.ObjectID, 0)
	}

	if event.Teams == nil {
		event.Teams = make([]primitive.ObjectID, 0)
	}
}
