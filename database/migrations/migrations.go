package migrations

import (
	"github.com/go-pg/migrations"
	"github.com/vovainside/vobook/database"
)

func Migrate() (int64, int64, error) {
	return migrations.Run(database.Conn(), "up")
}
