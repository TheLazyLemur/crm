package ops

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"simplecrm/internal/db"
	"simplecrm/internal/pubsub"
)

func CreateUser(
	ctx context.Context,
	dbc *sqlx.DB,
	querier db.Querier,
	firstName, lastName, email string,
	userCreatedEventService pubsub.UserCreatedEventServicer,
) (user db.User, err error) {
	tx, err := dbc.Beginx()
	if err != nil {
		return db.User{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	id := uuid.New().String()
	user, err = querier.InsertAndReturnUser(ctx, tx, db.InsertAndReturnUserParams{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	})
	if err != nil {
		return db.User{}, err
	}

	if err := userCreatedEventService.Publish(ctx, user); err != nil {
		return db.User{}, err
	}

	return user, nil
}
