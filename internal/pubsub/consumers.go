package pubsub

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"simplecrm/internal/db"
)

// TODO: Replace with a real event bus or Sqlite
type UserCreatedEventServicer interface {
	Consume(ctx context.Context, f func(UserCreatedEvent))
	Publish(ctx context.Context, user db.User) error
}

type userCreatedEventService chan UserCreatedEvent

type UserCreatedEvent struct {
	User db.User
}

func NewUserCreatedEventService() UserCreatedEventServicer {
	return make(userCreatedEventService, 100)
}

func (c userCreatedEventService) Consume(ctx context.Context, f func(UserCreatedEvent)) {
	for {
		select {
		case userCreatedEvent := <-c:
			slog.Info("User created event received", "user", userCreatedEvent.User)
			f(userCreatedEvent)
		default:
			time.Sleep(time.Second)
		}
	}
}

func (c userCreatedEventService) Publish(ctx context.Context, user db.User) error {
	select {
	case c <- UserCreatedEvent{User: user}:
		fmt.Println("User created event published")
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("timeout")
	}
}
