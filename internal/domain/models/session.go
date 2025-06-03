package models

import (
	"context"

	"github.com/google/uuid"

	"github.com/SlashLight/todo-list/pkg/my_err"
)

type Session struct {
	UserID uuid.UUID
	Email  string
}

type sessKey string

var SessionKey sessKey = "session"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, my_err.ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
