package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"

	"simplecrm/internal/db"
	"simplecrm/internal/ops"
)

type handlerFunc[T Validatable] func(w http.ResponseWriter, r *http.Request, params T)

// User handlers

func CreateUser(
	dbc *sqlx.DB,
	querier db.Querier,
) handlerFunc[createUserRequest] {
	return func(w http.ResponseWriter, r *http.Request, req createUserRequest) {
		if validationError := req.Validate(); len(validationError) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error(validationError.Error())
			return
		}

		user, err := ops.CreateUser(
			r.Context(),
			dbc,
			querier,
			req.FirstName,
			req.LastName,
			req.Email,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
			return
		}

		resp := mapUserToResponse(user)

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(err.Error())
		}
	}
}

func HandleUserCommand() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Notified"))
	}
}

func UpdateUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Updated"))
	}
}

func GetUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User"))
	}
}

// Task handlers

func CreateTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Created"))
	}
}

func UpdateTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Updated"))
	}
}

func HandleTaskCommand() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Started"))
	}
}

func GetTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Task"))
	}
}

// Lead handlers

func GetLead() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Lead"))
	}
}

func CreateLead() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Created"))
	}
}

func UpdateLead() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Updated"))
	}
}

func HandleLeadCommand() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Assigned"))
	}
}

// Contact handlers

func CreateContact() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Created"))
	}
}

func UpdateContact() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Updated"))
	}
}

func HandleContactCommand() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Contacted"))
	}
}

func GetContact() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Contact"))
	}
}
