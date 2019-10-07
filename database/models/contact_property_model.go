package models

import (
	"context"
	"time"

	contactproperty "github.com/vovainside/vobook/enum/contact_property"
)

type ContactProperty struct {
	ID        string               `json:"id"`
	ContactID string               `json:"contact_id"`
	Type      contactproperty.Type `json:"type"`
	Name      string               `json:"name"`
	Value     string               `json:"value"`
	Order     int                  `json:"order"`
	CreatedAt time.Time            `json:"created_at"`
	DeletedAt time.Time            `json:"deleted_at"`
}

func (m *ContactProperty) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	return ctx, nil
}
