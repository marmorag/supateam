package routes

import (
	"bytes"
	"encoding/json"
	"github.com/marmorag/supateam/internal/seeder"
	coretesting "github.com/marmorag/supateam/internal/testing"
	"net/http/httptest"
	"testing"
)

func Test_authUser(t *testing.T) {
	tests := []struct {
		name         string
		identity     string
		expectedCode int
	}{
		{
			name:         "test can login - admin user",
			identity:     "0600000001",
			expectedCode: 200,
		},
		{
			name:         "test can login - normal user",
			identity:     "0600000002",
			expectedCode: 200,
		},
		{
			name:         "test can login - user with no access",
			identity:     "0600000003",
			expectedCode: 200,
		},
	}

	app, router := coretesting.BuildTestApplication()
	AuthRouteHandler{}.Register(*router)

	_ = seeder.HttpTestSeeder{
		EmptyCollections: true,
	}.Seed()

	for _, tt := range tests {
		authBody := AuthRequest{
			Identity: tt.identity,
		}
		jsonBody, _ := json.Marshal(authBody)

		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(jsonBody))
		req.Header.Add("Content-Type", "application/json")

		resp, _ := app.Test(req, 1)
		if resp.StatusCode != tt.expectedCode {
			t.Errorf("HTTP Status differ : expected(%v) obtained(%v)", tt.expectedCode, resp.StatusCode)
		}
	}
}
