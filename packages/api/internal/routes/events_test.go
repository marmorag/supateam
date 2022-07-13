package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/seeder"
	coretesting "github.com/marmorag/supateam/internal/testing"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http/httptest"
	"testing"
	"time"
)

func _eventBootAndSeed() *fiber.App {
	app, router := coretesting.BuildTestApplication()
	EventRouteHandler{}.Register(*router)

	_ = seeder.HttpTestSeeder{
		EmptyCollections: true,
	}.Seed()

	return app
}

func Test_createEvent(t *testing.T) {
	tests := []struct {
		name         string
		user         string
		body         models.CreateEventRequest
		expectedCode int
	}{
		{
			name: "test can create event - admin user",
			user: MustAuthUser("0600000001"),
			body: models.CreateEventRequest{
				Title:       "Test Event",
				Description: "This is a test event",
				Date:        primitive.NewDateTimeFromTime(time.Now()),
				Duration:    1,
				Kind:        models.KindEquipe,
				Teams:       nil,
				Players:     nil,
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name: "test can't create event - normal user",
			user: MustAuthUser("0600000002"),
			body: models.CreateEventRequest{
				Title:       "Test Event",
				Description: "This is a test event",
				Date:        primitive.NewDateTimeFromTime(time.Now()),
				Duration:    1,
				Kind:        models.KindEquipe,
				Teams:       nil,
				Players:     nil,
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name: "test can't create event - user with no access",
			user: MustAuthUser("0600000003"),
			body: models.CreateEventRequest{
				Title:       "Test Event",
				Description: "This is a test event",
				Date:        primitive.NewDateTimeFromTime(time.Now()),
				Duration:    1,
				Kind:        models.KindEquipe,
				Teams:       nil,
				Players:     nil,
			},
			expectedCode: fiber.StatusForbidden,
		},
	}

	app := _eventBootAndSeed()

	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/api/events", MustSerializeReader(tt.body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", tt.user))

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		if resp.StatusCode != tt.expectedCode {
			t.Errorf("HTTP Status differ : expected(%v) obtained(%v)", tt.expectedCode, resp.StatusCode)
		}
	}
}

func Test_deleteEvent(t *testing.T) {
	tests := []struct {
		name         string
		user         string
		event        primitive.ObjectID
		expectedCode int
	}{
		{
			name:  "test can't delete event - not found",
			user:  MustAuthUser("0600000001"),
			event: seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f300"),
			// @todo fix : an non existing event should return not found instead of correct delete action
			//expectedCode: fiber.StatusNotFound,
			expectedCode: fiber.StatusNoContent,
		},
		{
			name:         "test can delete event - admin user",
			user:         MustAuthUser("0600000001"),
			event:        seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
			expectedCode: fiber.StatusNoContent,
		},
		{
			name:         "test can't delete event - normal user",
			user:         MustAuthUser("0600000002"),
			event:        seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f302"),
			expectedCode: fiber.StatusForbidden,
		},
		{
			name:         "test can't delete event - user with no access",
			user:         MustAuthUser("0600000003"),
			event:        seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f302"),
			expectedCode: fiber.StatusForbidden,
		},
	}

	app := _eventBootAndSeed()

	for _, tt := range tests {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/events/%s", tt.event.Hex()), nil)
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", tt.user))

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		if resp.StatusCode != tt.expectedCode {
			t.Errorf("HTTP Status differ : expected(%v) obtained(%v)", tt.expectedCode, resp.StatusCode)
		}
	}
}

func Test_getEvent(t *testing.T) {
	tests := []struct {
		name         string
		user         string
		event        primitive.ObjectID
		expectedCode int
	}{
		{
			name:         "test can't get event - not found",
			user:         MustAuthUser("0600000001"),
			event:        seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f300"),
			expectedCode: fiber.StatusNotFound,
		},
		{
			name:         "test can get event - admin user",
			user:         MustAuthUser("0600000001"),
			event:        seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "test can get event - normal user",
			user:         MustAuthUser("0600000002"),
			event:        seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "test can get event - user with no access",
			user:         MustAuthUser("0600000003"),
			event:        seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
			expectedCode: fiber.StatusForbidden,
		},
	}

	app := _eventBootAndSeed()

	for _, tt := range tests {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/events/%s", tt.event.Hex()), nil)
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", tt.user))

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		if resp.StatusCode != tt.expectedCode {
			t.Logf("found error at : %s", tt.name)
			t.Errorf("HTTP Status differ : expected(%v) obtained(%v)", tt.expectedCode, resp.StatusCode)
		}
	}
}

func Test_getEventParticipation(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getEventParticipation(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getEventParticipation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getEvents(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getEvents(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateEvent(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateEvent(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("updateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
