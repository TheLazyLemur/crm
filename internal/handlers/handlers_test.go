package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"simplecrm/database"
	"simplecrm/internal/db"
	"simplecrm/internal/pubsub/mocks"
)

func setupTest(t *testing.T) (*sqlx.DB, *chi.Mux, mocks.MockUserCreatedEventServicer, func()) {
	a := require.New(t)
	dbc, err := sqlx.Connect("sqlite3", ":memory:")
	a.NoError(err)

	_, err = dbc.Exec(string(database.Migrations))
	a.NoError(err)

	querier := &db.Queries{}
	r := chi.NewRouter()
	controller := gomock.NewController(t)
	eventService := mocks.NewMockUserCreatedEventServicer(controller)
	MountRoutes(r, dbc, querier, eventService)

	cleanup := func() {
		dbc.Close()
	}

	return dbc, r, *eventService, cleanup
}

func TestCreateUser(t *testing.T) {
	// Setup
	a := require.New(t)
	_, r, eventService, cleanup := setupTest(t)
	eventService.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	defer cleanup()

	// Test
	url := "/api/v1/user/create"
	pl := `{"first_name": "John", "last_name": "Doe", "email": "john.doe@example.com"}`
	req := httptest.NewRequest("POST", url, strings.NewReader(pl))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	a.Equal(w.Code, http.StatusCreated)

	body, err := io.ReadAll(w.Body)
	a.NoError(err)

	var user createUserResponse
	err = json.Unmarshal(body, &user)
	a.NoError(err)

	a.Equal("John", user.FirstName)
	a.Equal("Doe", user.LastName)
	a.Equal("john.doe@example.com", user.Email)
}

func TestCreateUser_BadRequest(t *testing.T) {
	// Setup
	a := require.New(t)
	_, r, _, cleanup := setupTest(t)
	defer cleanup()

	// Test
	url := "/api/v1/user/create"
	pl := ""
	req := httptest.NewRequest("POST", url, strings.NewReader(pl))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	a.Equal(w.Code, http.StatusBadRequest)
}

func TestCreateUser_DuplicateUserEmail(t *testing.T) {
	// Setup
	a := require.New(t)
	dbc, r, _, cleanup := setupTest(t)
	defer cleanup()

	// Test
	_, err := dbc.Exec(
		"INSERT INTO users (id, first_name, last_name, email) VALUES ('testid', 'John', 'Doe', 'john.doe@example.com')",
	)
	a.NoError(err)
	url := "/api/v1/user/create"
	pl := `{"first_name": "John", "last_name": "Doe", "email": "john.doe@example.com"}`
	req := httptest.NewRequest("POST", url, strings.NewReader(pl))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	a.Equal(w.Code, http.StatusInternalServerError)
}

func TestCreateUser_FailValidation(t *testing.T) {
	// Setup
	a := require.New(t)
	_, r, _, cleanup := setupTest(t)
	defer cleanup()

	// Test
	url := "/api/v1/user/create"

	tcs := []struct {
		name        string
		pl          string
		contentType string
	}{
		{
			name:        "Missing first name",
			pl:          `{"last_name": "Doe", "email": "john.doe@example.com"}`,
			contentType: "application/json",
		},
		{
			name:        "Missing last name",
			pl:          `{"first_name": "John", "email": "john.doe@example.com"}`,
			contentType: "application/json",
		},
		{
			name:        "Missing email",
			pl:          `{"first_name": "John", "last_name": "Doe"}`,
			contentType: "application/json",
		},
		{
			name:        "Invalid email",
			pl:          `{"first_name": "John", "last_name": "Doe", "email": "john.doe"}`,
			contentType: "application/json",
		},
		{
			name: "Missing content type",
			pl:   `{"first_name": "John","last_name": "Doe", "email": "john.doe@example.com"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", url, strings.NewReader(tc.pl))
			req.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			a.Equal(http.StatusBadRequest, w.Code)
		})
	}
}

func TestGetUser(t *testing.T) {
	tcs := []struct {
		name              string
		id                string
		expetedStatusCode int
		expected          getUserResponse
	}{
		{
			name:              "Missing id",
			id:                "",
			expetedStatusCode: http.StatusBadRequest,
			expected:          getUserResponse{},
		},
		{
			name:              "Invalid id",
			id:                "invalid",
			expetedStatusCode: http.StatusInternalServerError,
			expected:          getUserResponse{},
		},
		{
			name:              "Valid id",
			id:                "testid",
			expetedStatusCode: http.StatusOK,
			expected: getUserResponse{
				ID:        "testid",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
		},
	}
	// Setup
	a := require.New(t)
	dbc, r, _, cleanup := setupTest(t)
	defer cleanup()
	_, err := dbc.Exec(
		"INSERT INTO users (id, first_name, last_name, email) VALUES ('testid', 'John', 'Doe', 'john.doe@example.com')",
	)
	a.NoError(err)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Test
			url := "/api/v1/query/user?id=" + tc.id
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			a.Equal(w.Code, tc.expetedStatusCode)

			body, err := io.ReadAll(w.Body)
			a.NoError(err)

			if tc.expetedStatusCode == http.StatusOK {
				var user getUserResponse
				err = json.Unmarshal(body, &user)
				a.NoError(err)

				a.Equal(tc.expected.ID, user.ID)
				a.Equal(tc.expected.FirstName, user.FirstName)
				a.Equal(tc.expected.LastName, user.LastName)
				a.Equal(tc.expected.Email, user.Email)
			}
		})
	}
}
