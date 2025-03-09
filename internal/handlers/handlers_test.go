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

	"simplecrm/database"
	"simplecrm/internal/db"
)

func setupTest(t *testing.T) (*sqlx.DB, *chi.Mux, func()) {
	a := require.New(t)
	dbc, err := sqlx.Connect("sqlite3", ":memory:")
	a.NoError(err)

	_, err = dbc.Exec(string(database.Migrations))
	a.NoError(err)

	querier := &db.Queries{}
	r := chi.NewRouter()
	MountRoutes(r, dbc, querier)

	cleanup := func() {
		dbc.Close()
	}

	return dbc, r, cleanup
}

func TestCreateUser(t *testing.T) {
	// Setup
	a := require.New(t)
	_, r, cleanup := setupTest(t)
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
	_, r, cleanup := setupTest(t)
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
	dbc, r, cleanup := setupTest(t)
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
	_, r, cleanup := setupTest(t)
	defer cleanup()

	// Test
	url := "/api/v1/user/create"

	tcs := []struct {
		name string
		pl   string
	}{
		{
			name: "Missing first name",
			pl:   `{"last_name": "Doe", "email": "john.doe@example.com"}`,
		},
		{
			name: "Missing last name",
			pl:   `{"first_name": "John", "email": "john.doe@example.com"}`,
		},
		{
			name: "Missing email",
			pl:   `{"first_name": "John", "last_name": "Doe"}`,
		},
		{
			name: "Invalid email",
			pl:   `{"first_name": "John", "last_name": "Doe", "email": "john.doe"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", url, strings.NewReader(tc.pl))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			a.Equal(w.Code, http.StatusBadRequest)
		})
	}
}
