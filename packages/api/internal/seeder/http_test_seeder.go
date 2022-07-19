package seeder

import (
	"context"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/middleware/auth"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type HttpTestSeeder struct {
	EmptyCollections bool
}

func NewHttpTestSeeder(emptyCollections bool) *HttpTestSeeder {
	return &HttpTestSeeder{EmptyCollections: emptyCollections}
}

func (h HttpTestSeeder) Seed() error {
	config := internal.GetConfig()
	config.TracingEnabled = false
	internal.Set(config)

	// empty collection
	if h.EmptyCollections {
		_ = repository.NewEventRepository("").Collection.Drop(context.Background())
		_ = repository.NewUserRepository("").Collection.Drop(context.Background())
		_ = repository.NewTeamRepository("").Collection.Drop(context.Background())
		_ = repository.NewParticipationRepository("").Collection.Drop(context.Background())
	}

	err := httpSeedUsers()
	if err != nil {
		return err
	}

	err = httpSeedTeams()
	if err != nil {
		return err
	}

	err = httpSeedEvents()
	if err != nil {
		return err
	}

	return nil
}

func httpSeedTeams() error {
	tr := repository.NewTeamRepository("")

	team := models.Team{
		Id:   MustObjectIdFromHex("62cecdffa29e0c7df4c0f201"),
		Name: "Equipe 1",
		Members: []primitive.ObjectID{
			MustObjectIdFromHex("62cecdffa29e0c7df4c0f101"),
			MustObjectIdFromHex("62cecdffa29e0c7df4c0f102"),
			MustObjectIdFromHex("62cecdffa29e0c7df4c0f103"),
		},
	}

	_, err := tr.Collection.InsertOne(context.TODO(), &team)
	if err != nil {
		return err
	}

	return nil
}

func httpSeedEvents() error {
	events := make([]models.Event, 0)
	teams, _ := repository.NewTeamRepository("").FindAll()
	teamId := teams[0].Id

	events = append(events, models.Event{
		Id:          MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
		Title:       "Event test 1",
		Description: "First test event",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC)),
		Duration:    1,
		Kind:        models.KindEntrainement,
		Teams:       []primitive.ObjectID{teamId},
		Players:     make([]primitive.ObjectID, 0),
	})

	events = append(events, models.Event{
		Id:          MustObjectIdFromHex("62cecdffa29e0c7df4c0f302"),
		Title:       "Event test 2",
		Description: "Second test event",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.January, 15, 0, 0, 0, 0, time.UTC)),
		Duration:    1,
		Kind:        models.KindEntrainement,
		Teams:       []primitive.ObjectID{teamId},
		Players:     make([]primitive.ObjectID, 0),
	})

	events = append(events, models.Event{
		Id:          MustObjectIdFromHex("62cecdffa29e0c7df4c0f303"),
		Title:       "Event test 3",
		Description: "Third test event",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.February, 6, 0, 0, 0, 0, time.UTC)),
		Duration:    1,
		Kind:        models.KindEntrainement,
		Teams:       []primitive.ObjectID{teamId},
		Players:     make([]primitive.ObjectID, 0),
	})

	er := repository.NewEventRepository("")
	pr := repository.NewParticipationRepository("")
	for _, event := range events {
		_, err := er.Collection.InsertOne(context.TODO(), &event)
		if err != nil {
			return err
		}

		err = pr.SyncParticipation(&event)
		if err != nil {
			return err
		}
	}

	return nil
}

func httpSeedUsers() error {
	users := make([]models.User, 0)

	users = append(users, models.User{
		Id:       MustObjectIdFromHex("62cecdffa29e0c7df4c0f101"),
		Identity: "0600000001",
		Name:     "User:Admin",
		Authorization: map[auth.ApiGroups][]auth.ApiAction{
			auth.EventsApiGroup:         {auth.AllAction},
			auth.UsersApiGroup:          {auth.AllAction},
			auth.ParticipationsApiGroup: {auth.AllAction},
			auth.TeamsApiGroup:          {auth.AllAction},
		},
	})

	users = append(users, models.User{
		Id:       MustObjectIdFromHex("62cecdffa29e0c7df4c0f102"),
		Identity: "0600000002",
		Name:     "User:Normal",
		Authorization: map[auth.ApiGroups][]auth.ApiAction{
			auth.EventsApiGroup:         {auth.ReadAction, auth.WriteAction, auth.UpdateAction},
			auth.UsersApiGroup:          {auth.ReadAction, auth.UpdateSelfAction},
			auth.ParticipationsApiGroup: {auth.ReadAction, auth.UpdateSelfAction, auth.WriteSelfAction},
			auth.TeamsApiGroup:          {auth.ReadAction},
		},
	})

	users = append(users, models.User{
		Id:       MustObjectIdFromHex("62cecdffa29e0c7df4c0f103"),
		Identity: "0600000003",
		Name:     "User:None",
		Authorization: map[auth.ApiGroups][]auth.ApiAction{
			auth.EventsApiGroup:         {},
			auth.UsersApiGroup:          {},
			auth.ParticipationsApiGroup: {},
			auth.TeamsApiGroup:          {},
		},
	})

	ur := repository.NewUserRepository("")
	for _, user := range users {
		_, err := ur.Collection.InsertOne(context.TODO(), &user)
		if err != nil {
			return err
		}
	}

	return nil
}
