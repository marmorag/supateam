package routes

import (
	"github.com/marmorag/supateam/internal/seeder"
	coretesting "github.com/marmorag/supateam/internal/testing"
	"net/http/httptest"
	"testing"
)

func Test_authUser(t *testing.T) {
	tests := []struct {
		name         string
		body         AuthRequest
		expectedCode int
	}{
		{
			name:         "test can login - admin user",
			body:         AuthRequest{"0600000001"},
			expectedCode: 200,
		},
		{
			name:         "test can login - normal user",
			body:         AuthRequest{"0600000002"},
			expectedCode: 200,
		},
		{
			name:         "test can login - user with no access",
			body:         AuthRequest{"0600000003"},
			expectedCode: 200,
		},
		{
			name:         "test can't login - user not known",
			body:         AuthRequest{"0600000000"},
			expectedCode: 404,
		},
		{
			name:         "test can't login - Identity not provided",
			body:         AuthRequest{},
			expectedCode: 400,
		},
	}

	app, router := coretesting.BuildTestApplication()
	AuthRouteHandler{}.Register(*router)

	_ = seeder.HttpTestSeeder{
		EmptyCollections: true,
	}.Seed()

	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/api/auth/login", MustSerializeReader(tt.body))
		req.Header.Add("Content-Type", "application/json")

		resp, _ := app.Test(req, -1)
		if resp.StatusCode != tt.expectedCode {
			t.Errorf("HTTP Status differ : expected(%v) obtained(%v)", tt.expectedCode, resp.StatusCode)
		}
	}
}
