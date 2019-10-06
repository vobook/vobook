package models

import (
	"context"
	"time"

	"github.com/vovainside/vobook/config"
)

type PasswordReset struct {
	ID        string
	UserID    string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (m *PasswordReset) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.ExpiresAt.IsZero() {
		m.ExpiresAt = m.CreatedAt.Add(config.Get().PasswordResetLifetime)
	}
	return ctx, nil
}
