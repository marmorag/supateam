package seeder

import (
	"context"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/middleware/auth"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
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

	return nil
}

func httpSeedUsers() error {
	users := make([]models.User, 0)

	users = append(users, models.User{
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
		_, err := ur.Create(&user)
		if err != nil {
			return err
		}
	}

	return nil
}
