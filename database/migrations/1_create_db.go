package migrations

import (
	"github.com/go-pg/migrations/v7"
)

func init() {
	mg := `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id             UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    first_name     TEXT,
    last_name      TEXT,
    email          TEXT UNIQUE      NOT NULL,
    email_verified BOOL                      DEFAULT FALSE,
    password       TEXT             NOT NULL,
    telegram_id    INT,
    created_at     TIMESTAMPTZ      NOT NULL,
    deleted_at     TIMESTAMPTZ
);

CREATE TABLE email_verifications
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    token      TEXT             NOT NULL,
    user_id    UUID             NOT NULL REFERENCES users (id),
    email      TEXT,
    created_at TIMESTAMPTZ      NOT NULL
);

CREATE TABLE auth_tokens
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL REFERENCES users (id),
    client_id  INT         NOT NULL,
    client_ip  TEXT,
    user_agent TEXT,
    token      TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE password_resets
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL REFERENCES users (id),
    token      TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE contacts
(
    id          UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id     UUID             NOT NULL REFERENCES users (id),
    name        TEXT,
    first_name  TEXT,
    last_name   TEXT,
    middle_name TEXT,
    birthday    TIMESTAMPTZ,
    gender      INT,
    created_at  TIMESTAMPTZ      NOT NULL,
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE contact_properties
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    contact_id UUID             NOT NULL REFERENCES contacts (id),
    name       TEXT             NOT NULL,
    type       INT              NOT NULL,
    value      TEXT             NOT NULL,
    "order"    INT              NOT NULL,
    created_at TIMESTAMPTZ      NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE birthday_notification_logs
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    contact_id UUID             NOT NULL REFERENCES contacts (id),
    created_at TIMESTAMPTZ      NOT NULL
);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
