package migrations

import (
	"github.com/go-pg/migrations/v7"
)

func init() {
	up := `
CREATE TABLE files
(
    id           UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id      UUID             NOT NULL REFERENCES users (id),
    name         TEXT,
    description  TEXT,
    size         BIGINT,
    hash         TEXT             NOT NULL,
    type         INT              NOT NULL,
    path         TEXT             NOT NULL,
    preview_path TEXT,
    created_at   TIMESTAMPTZ      NOT NULL,
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ
);

ALTER TABLE contacts ADD photo_id UUID REFERENCES files (id);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(up)
		return
	})
}
