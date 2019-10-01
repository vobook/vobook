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
    email text,
    password text,
    created_at timestamptz not null
);
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
