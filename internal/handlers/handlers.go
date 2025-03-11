package handlers

import (
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"

	"simplecrm/internal/db"
	"simplecrm/internal/ops"
	"simplecrm/internal/pubsub"
)

func CreateUser(
	dbc *sqlx.DB,
	querier db.Querier,
	userCreatedEventService pubsub.UserCreatedEventServicer,
) handlerFunc[createUserRequest, createUserResponse] {
	return func(w http.ResponseWriter, r *http.Request, req createUserRequest) (*httpResponse[createUserResponse], *httpError) {
		if validationError := req.Validate(); len(validationError) > 0 {
			return nil, &httpError{
				Message:    validationError.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}

		user, err := ops.CreateUser(
			r.Context(),
			dbc,
			querier,
			req.FirstName,
			req.LastName,
			req.Email,
			userCreatedEventService,
		)
		if err != nil {
			slog.Error(err.Error())
			return nil, &httpError{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		return &httpResponse[createUserResponse]{
			Data:       mapUserToResponse(user),
			StatusCode: http.StatusCreated,
		}, nil
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
