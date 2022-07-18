package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/seeder"
	coretesting "github.com/marmorag/supateam/internal/testing"
	"net/http/httptest"
	"testing"
)

func _participationBootAndSeed() *fiber.App {
	app, router := coretesting.BuildTestApplication()
	ParticipationRouteHandler{}.Register(*router)

	_ = seeder.HttpTestSeeder{
		EmptyCollections: true,
	}.Seed()

	return app
}

func Test_createParticipation(t *testing.T) {
	tests := []struct {
		name         string
		user         string
		body         models.CreateParticipationRequest
		expectedCode int
	}{
		{
			name: "test can create participation - admin user",
			user: MustAuthUser("0600000001"),
			body: models.CreateParticipationRequest{
				Event:  seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
				Player: seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f101"),
				Status: models.ParticipationAccepted,
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name: "test can create participation for any one - admin user",
			user: MustAuthUser("0600000001"),
			body: models.CreateParticipationRequest{
				Event:  seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f302"),
				Player: seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f102"),
				Status: models.ParticipationAccepted,
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name: "test can create participation - normal user",
			user: MustAuthUser("0600000002"),
			body: models.CreateParticipationRequest{
				Event:  seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
				Player: seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f102"),
				Status: models.ParticipationAccepted,
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name: "test can't create participation - normal user trying to create on another account",
			user: MustAuthUser("0600000002"),
			body: models.CreateParticipationRequest{
				Event:  seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
				Player: seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f103"),
				Status: models.ParticipationAccepted,
			},
			expectedCode: fiber.StatusForbidden,
		},
		{
			name: "test can't create participation - user with no access",
			user: MustAuthUser("0600000003"),
			body: models.CreateParticipationRequest{
				Event:  seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f301"),
				Player: seeder.MustObjectIdFromHex("62cecdffa29e0c7df4c0f103"),
				Status: models.ParticipationAccepted,
			},
			expectedCode: fiber.StatusForbidden,
		},
	}

	app := _participationBootAndSeed()

	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/api/participations", MustSerializeReader(tt.body))
		req.Header.Add("Content-Type", "application/json")
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

func Test_deleteParticipation(t *testing.T) {
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
			if err := deleteParticipation(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("deleteParticipation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getParticipation(t *testing.T) {
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
			if err := getParticipation(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getParticipation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getParticipations(t *testing.T) {
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
			if err := getParticipations(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getParticipations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateParticipation(t *testing.T) {
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
			if err := updateParticipation(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("updateParticipation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
