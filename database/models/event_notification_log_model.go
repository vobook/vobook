package models

import (
	"context"
	"time"
)

type BirthdayNotificationLog struct {
	ID        string
	ContactID string
	CreatedAt time.Time
}

func (m *BirthdayNotificationLog) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	return ctx, nil
}
