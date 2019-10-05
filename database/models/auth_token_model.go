package models

import (
	"context"
	"time"

	"github.com/vovainside/vobook/config"
)

type ClientID int

const (
	_ ClientID = iota
	VueClient
)

var Clients = []ClientID{
	VueClient,
}

type AuthToken struct {
	ID        string
	UserID    string
	User      *User
	ClientID  ClientID
	ClientIP  string
	UserAgent string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (m *AuthToken) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.ExpiresAt.IsZero() {
		m.ExpiresAt = m.CreatedAt.Add(config.Get().AuthTokenLifetime)
	}
	return ctx, nil
}
