package models

import (
	"context"
	"time"
	"vobook/services/fs"
)

type File struct {
	ID          string      `json:"id"`
	UserID      string      `json:"-"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Size        int64       `json:"size"`
	Hash        string      `json:"-"` // should not be exposed
	Type        fs.FileType `json:"type"`
	Path        string      `json:"-"` // should not be exposed
	PreviewPath string      `json:"-"` // should not be exposed
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at"`
	DeletedAt   *time.Time  `json:"deleted_at"`

	Filename string `json:"-" sql:"-"`
	Bytes    []byte `json:"-" sql:"-"`
	Base64   string `json:"base64" sql:"-"`
}

func (m *File) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (m *File) BeforeUpdate(ctx context.Context) (context.Context, error) {
	if m.UpdatedAt == nil || m.UpdatedAt.IsZero() {
		now := time.Now()
		m.UpdatedAt = &now
	}
	return ctx, nil
}
