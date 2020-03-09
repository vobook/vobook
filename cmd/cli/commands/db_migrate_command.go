package commands

import (
	"fmt"

	"vobook/database"
	"vobook/database/migrations"
)

func init() {
	add("db-migrate", "migrate", "mg", command{
		handler: migrateDB,
		help:    `Migrate database`,
	})
}

func migrateDB(args ...string) (err error) {
	db := database.Conn()
	n, err := db.Model().
		Table("pg_tables").
		Where("schemaname = 'public'").
		Where("tablename = 'gopg_migrations'").
		Count()
	if err != nil {
		return
	}
	if n == 0 {
		_, err = db.Exec(`
		CREATE SEQUENCE gopg_migrations_id_seq;
		CREATE TABLE gopg_migrations
		(
			id         INT DEFAULT nextval('gopg_migrations_id_seq'),
			version    INT,
			created_at TIMESTAMPTZ
		);`)
		if err != nil {
			return
		}
	}

	oldV, newV, err := migrations.Migrate()
	if err != nil {
		return
	}

	if oldV == newV {
		println("Nothing to migrate")
		return
	}

	println(fmt.Sprintf("Migrated from %d to  %d", oldV, newV))
	return
}
