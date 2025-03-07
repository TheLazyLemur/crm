package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"simplecrm/database"
	"simplecrm/internal/db"
)

func setupTest(t *testing.T) (*sql.DB, *chi.Mux, func()) {
	a := require.New(t)
	dbc, err := sql.Open("sqlite3", ":memory:")
	a.NoError(err)

	_, err = dbc.Exec(string(database.Migrations))
	a.NoError(err)

	querier := db.New()
	r := chi.NewRouter()
	mountRoutes(r, dbc, querier)

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
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	a.Equal(w.Code, http.StatusCreated)

	body, err := io.ReadAll(w.Body)
	a.NoError(err)

	var user db.User
	err = json.Unmarshal(body, &user)
	a.NoError(err)

	a.Equal(user.FirstName, "John")
	a.Equal(user.LastName, "Doe")
	a.Equal(user.Email, "john.doe@example.com")
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
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	a.Equal(w.Code, http.StatusInternalServerError)
}
