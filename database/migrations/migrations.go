package migrations

import (
	"vobook/database"

	"github.com/go-pg/migrations/v7"
)

func Migrate() (int64, int64, error) {
	return migrations.Run(database.Conn(), "up")
}
