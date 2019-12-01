package models

import (
	"context"
	"time"

	"github.com/vovainside/vobook/enum/gender"
)

type Contact struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Name       string     `json:"name"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	MiddleName string     `json:"middle_name"`
	DOBYear    int        `json:"dob_year"`
	DOBMonth   time.Month `json:"dob_month"`
	DOBDay     int        `json:"dob_day"`

	Gender    gender.Type `json:"gender"`
	CreatedAt time.Time   `json:"created_at"`
	DeletedAt *time.Time  `json:"deleted_at"`

	// Relations
	User  *User             `json:"user"`
	Props []ContactProperty `json:"props"`
}

func (m *Contact) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	return ctx, nil
}

func (m *Contact) AfterSelect(ctx context.Context) error {
	if m.Name == "" {
		m.Name = m.FirstName + " " + m.LastName
	}

	return nil
}
