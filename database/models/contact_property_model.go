package models

import (
	"context"
	"time"

	"github.com/vovainside/vobook/enum/contact_property"
)

type ContactProperty struct {
	ID        string               `json:"id"`
	ContactID string               `json:"contact_id"`
	Type      contactproperty.Type `json:"type" pg:",use_zero"`
	Name      string               `json:"name" pg:",use_zero"`
	Value     string               `json:"value" pg:",use_zero"`
	Order     int                  `json:"order" pg:",use_zero"`
	CreatedAt time.Time            `json:"created_at"`
	DeletedAt *time.Time           `json:"deleted_at"`

	Contact *Contact `json:"-"`
}

func (m *ContactProperty) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	return ctx, nil
}
