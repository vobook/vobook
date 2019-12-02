package migrations

import (
	"github.com/go-pg/migrations"
)

func init() {
	mg := `
ALTER TABLE contacts ADD COLUMN dob_year INT,
ADD COLUMN dob_month SMALLINT,
ADD COLUMN dob_day SMALLINT;

UPDATE contacts SET dob_year=date_part('year', birthday);
UPDATE contacts SET dob_month=date_part('month', birthday);
UPDATE contacts SET dob_day=date_part('day', birthday);

ALTER TABLE contacts DROP birthday;
`

	migrations.MustRegister(func(db migrations.DB) (err error) {
		_, err = db.Exec(mg)
		return
	})
}
