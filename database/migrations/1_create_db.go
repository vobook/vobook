package migrations

import (
	"github.com/go-pg/migrations"
)

func init() {
	mg := `
create table users (
    id uuid primary key not null default gen_random_uuid(),
    first_name text,
    last_name text,
    email text unique not null,
	email_verified bool default false,
    password text not null,
    created_at timestamptz not null
);

create table email_verifications (
    id uuid primary key not null default gen_random_uuid(),
    user_id uuid references users (id),
    email text not null,
    created_at timestamptz not null
);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
