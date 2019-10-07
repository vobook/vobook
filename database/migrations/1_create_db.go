package migrations

import (
	"github.com/go-pg/migrations"
)

func init() {
	mg := `
create table users
(
    id             uuid primary key not null default gen_random_uuid(),
    first_name     text,
    last_name      text,
    email          text unique      not null,
    email_verified bool                      default false,
    password       text             not null,
    created_at     timestamptz      not null,
    deleted_at     timestamptz
);

create table email_verifications
(
    id         uuid primary key not null default gen_random_uuid(),
    token      text             not null,
    user_id    uuid             not null references users (id),
    email      text,
    created_at timestamptz      not null
);

create table auth_tokens
(
    id         uuid primary key default gen_random_uuid(),
    user_id    uuid        not null references users (id),
    client_id  int         not null,
    client_ip  text,
    user_agent text,
    token      text,
    created_at timestamptz not null,
    expires_at timestamptz not null
);

create table password_resets
(
    id         uuid primary key default gen_random_uuid(),
    user_id    uuid        not null references users (id),
    token      text        not null,
    created_at timestamptz not null,
    expires_at timestamptz not null
);

create table contacts
(
    id          uuid primary key not null default gen_random_uuid(),
    user_id     uuid             not null references users (id),
    name        text,
    first_name  text,
    last_name   text,
    middle_name text,
    birthday    timestamptz,
    created_at  timestamptz      not null,
    deleted_at  timestamptz
);

create table contact_properties
(
    id         uuid primary key not null default gen_random_uuid(),
    contact_id uuid             not null references contacts (id),
    name       text,
    type       int,
    value      text,
    "order"    int,
    created_at timestamptz      not null,
    deleted_at timestamptz
);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
