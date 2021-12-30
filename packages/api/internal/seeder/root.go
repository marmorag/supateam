package seeder

import (
	"context"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Seeder struct{}

func (s Seeder) Seed() error {
	// empty collection
	_ = repository.NewEventRepository().Collection.Drop(context.Background())
	_ = repository.NewUserRepository().Collection.Drop(context.Background())
	_ = repository.NewTeamRepository().Collection.Drop(context.Background())
	_ = repository.NewParticipationRepository().Collection.Drop(context.Background())

	// seed all entities
	err := seedUsers()
	if err != nil {
		return err
	}

	err = seedTeams()
	if err != nil {
		return err
	}

	err = seedEvents()
	if err != nil {
		return err
	}

	err = seedAdmins()
	if err != nil {
		return err
	}

	return nil
}

func seedUsers() error {
	users := make([]models.User, 0)
	dec, err := ReadSecureData("./data.yaml.enc")
	if err != nil {
		return err
	}

	parsed, err := ParseYamlFromString(dec)
	if err != nil {
		return err
	}

	users = append(users, models.User{
		Identity:      parsed.Phones["Guillaume"],
		Name:          "Guillaume",
		Authorization: map[string][]string{"events": {"*"}, "users": {"*"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Alex"],
		Name:          "Alex",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Baptiste"],
		Name:          "Baptiste",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Vincent"],
		Name:          "Vincent",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Adrien"],
		Name:          "Adrien",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Jérémy"],
		Name:          "Jérémy",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Clément"],
		Name:          "Clément",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Mike"],
		Name:          "Mike",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Arthur"],
		Name:          "Arthur",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Flo"],
		Name:          "Flo",
		Authorization: map[string][]string{"events": {"read", "write"}, "users": {"read"}},
	})

	ur := repository.NewUserRepository()
	for _, user := range users {
		_, err := ur.Create(&user)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedTeams() error {
	ur := repository.NewUserRepository()
	tr := repository.NewTeamRepository()

	users, err := ur.FindAll()
	members := make([]primitive.ObjectID, 0)

	for _, user := range users {
		members = append(members, user.Id)
	}

	team := models.Team{
		Name:    "Equipe 1",
		Members: members,
	}

	_, err = tr.Create(&team)
	if err != nil {
		return err
	}

	return nil
}

func seedEvents() error {
	events := make([]models.Event, 0)
	teams, _ := repository.NewTeamRepository().FindAll()
	teamId := teams[0].Id

	events = append(events, models.Event{
		Title:       "Entrainement 1/4",
		Description: "Premier entrainement au Superflu.\nRDV 9h",
		Date:        primitive.NewDateTimeFromTime(time.Date(2021, time.December, 19, 0, 0, 0, 0, time.UTC)),
		Duration:    1,
		Kind:        models.KindEntrainement,
		Teams:       []primitive.ObjectID{teamId},
	})

	events = append(events, models.Event{
		Title:       "Entrainement 2/4",
		Description: "Deuxième entrainement au Superflu.\nAprès la Ringer.",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.January, 16, 0, 0, 0, 0, time.UTC)),
		Duration:    1,
		Kind:        models.KindEntrainement,
		Teams:       []primitive.ObjectID{teamId},
	})

	events = append(events, models.Event{
		Title:       "Entrainement 3/4",
		Description: "Troisième entrainement au Superflu.",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.February, 6, 0, 0, 0, 0, time.UTC)),
		Duration:    1,
		Kind:        models.KindEntrainement,
		Teams:       []primitive.ObjectID{teamId},
	})

	events = append(events, models.Event{
		Title:       "Entrainement 4/4",
		Description: "Dernier entrainement au Superflu.",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 6, 0, 0, 0, 0, time.UTC)),
		Duration:    1,
		Kind:        models.KindEntrainement,
		Teams:       []primitive.ObjectID{teamId},
	})

	events = append(events, models.Event{
		Title:       "Grand Prix de la Drôme",
		Description: "Grand Prix en équipe, Ben nous accompagnera.",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.April, 2, 0, 0, 0, 0, time.UTC)),
		Duration:    2,
		Kind:        models.KindGrandPrix,
		Teams:       []primitive.ObjectID{teamId},
	})

	events = append(events, models.Event{
		Title:       "Qualif. Promo",
		Description: "Golf de Vichy Montpensier.",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.April, 22, 0, 0, 0, 0, time.UTC)),
		Duration:    3,
		Kind:        models.KindEquipe,
		Teams:       []primitive.ObjectID{teamId},
	})

	events = append(events, models.Event{
		Title:       "Promotion",
		Description: "Golf de Vichy Montpensier.",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.May, 6, 0, 0, 0, 0, time.UTC)),
		Duration:    3,
		Kind:        models.KindEquipe,
		Teams:       []primitive.ObjectID{teamId},
	})

	events = append(events, models.Event{
		Title:       "Grand Prix du Clou",
		Description: "",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.May, 21, 0, 0, 0, 0, time.UTC)),
		Duration:    2,
		Kind:        models.KindGrandPrix,
	})

	events = append(events, models.Event{
		Title:       "Grand Prix du Lac d'Annecy",
		Description: "",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.June, 11, 0, 0, 0, 0, time.UTC)),
		Duration:    2,
		Kind:        models.KindGrandPrix,
	})

	events = append(events, models.Event{
		Title:       "Grand Prix de la Sorelle",
		Description: "",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.July, 2, 0, 0, 0, 0, time.UTC)),
		Duration:    2,
		Kind:        models.KindGrandPrix,
	})

	events = append(events, models.Event{
		Title:       "Grand Prix du Forez",
		Description: "",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.August, 19, 0, 0, 0, 0, time.UTC)),
		Duration:    3,
		Kind:        models.KindGrandPrix,
	})

	events = append(events, models.Event{
		Title:       "Grand Prix de Vichy",
		Description: "",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.August, 26, 0, 0, 0, 0, time.UTC)),
		Duration:    3,
		Kind:        models.KindGrandPrix,
	})

	events = append(events, models.Event{
		Title:       "Grand Prix des Volcans",
		Description: "",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.September, 9, 0, 0, 0, 0, time.UTC)),
		Duration:    3,
		Kind:        models.KindGrandPrix,
	})

	events = append(events, models.Event{
		Title:       "Championnat des golfs 9 trous",
		Description: "A la casa.",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.September, 17, 0, 0, 0, 0, time.UTC)),
		Duration:    2,
		Kind:        models.KindEquipe,
		Teams:       []primitive.ObjectID{teamId},
	})

	events = append(events, models.Event{
		Title:       "Grand Prix de la Ligue AURA",
		Description: "Golf du Gouverneur.",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.September, 23, 0, 0, 0, 0, time.UTC)),
		Duration:    3,
		Kind:        models.KindGrandPrix,
	})

	events = append(events, models.Event{
		Title:       "Grand Prix du Rhône",
		Description: "Golf de Lyon Chassieu",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.October, 1, 0, 0, 0, 0, time.UTC)),
		Duration:    2,
		Kind:        models.KindGrandPrix,
	})

	events = append(events, models.Event{
		Title:       "2e Division Régionale",
		Description: "Golf du Chambon-sur-Lignon",
		Date:        primitive.NewDateTimeFromTime(time.Date(2022, time.October, 8, 0, 0, 0, 0, time.UTC)),
		Duration:    2,
		Kind:        models.KindEquipe,
		Teams:       []primitive.ObjectID{teamId},
	})

	er := repository.NewEventRepository()
	pr := repository.NewParticipationRepository()
	for _, event := range events {
		created, err := er.Create(&event)
		if err != nil {
			return err
		}

		err = pr.SyncParticipation(created)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedAdmins() error {
	users := make([]models.User, 0)
	dec, err := ReadSecureData("./data.yaml.enc")
	if err != nil {
		return err
	}

	parsed, err := ParseYamlFromString(dec)
	if err != nil {
		return err
	}

	users = append(users, models.User{
		Identity:      parsed.Phones["Cyril"],
		Name:          "Cyril",
		Authorization: map[string][]string{"events": {"*"}, "users": {"*"}},
	})

	users = append(users, models.User{
		Identity:      parsed.Phones["Benoit"],
		Name:          "Benoit",
		Authorization: map[string][]string{"events": {"*"}, "users": {"*"}},
	})

	ur := repository.NewUserRepository()
	for _, user := range users {
		_, err := ur.Create(&user)
		if err != nil {
			return err
		}
	}

	return nil
}
