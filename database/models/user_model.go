package models

import (
	"context"
	"time"
)

type User struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Password      string    `json:"-"`
	TelegramID    int       `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	DeletedAt     time.Time `json:"deleted_at"`

	// Composite fields
	HasTelegram bool `json:"has_telegram" pg:"-"`
}

func (m *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (m *User) AfterSelect(ctx context.Context) error {
	m.HasTelegram = m.TelegramID != 0

	return nil
}
