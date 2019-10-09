package database

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	_ "github.com/lib/pq"

	"github.com/vovainside/vobook/config"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

var (
	conn *pg.DB
	DB   orm.DB
)

func connect() *pg.DB {
	conf := config.Get()
	conn = pg.Connect(&pg.Options{
		Addr:            conf.DB.Addr,
		User:            conf.DB.User,
		Password:        conf.DB.Password,
		Database:        conf.DB.Name,
		ApplicationName: conf.App.Name + "." + conf.App.Env,
	})

	conn.AddQueryHook(dbLogger{})
	DB = conn

	return conn
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) (err error) {
	query, err := q.FormattedQuery()
	if err != nil {
		return
	}
	println(query + "\n")
	return nil
}

func Conn() *pg.DB {
	if conn == nil {
		connect()
	}

	return conn
}

func ORM() orm.DB {
	if conn == nil {
		connect()
	}

	return DB
}

func SetDB(db orm.DB) {
	DB = db
}

func RunInTransaction(txFunc func(tx *pg.Tx) error) (err error) {
	// tests already running in transaction
	// wrapping in another one will result errors
	// because data from outer transaction is not done (and never will)
	if config.IsTestEnv() {
		return txFunc(ORM().(*pg.Tx))
	}

	tx, err := Conn().Begin()
	if err != nil {
		return
	}

	err = txFunc(tx)
	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

func Truncate(table string) error {
	_, err := ORM().Exec("TRUNCATE TABLE ? RESTART IDENTITY CASCADE", pg.F(table))
	return err
}

func Insert(v interface{}) error {
	return DB.Insert(v)
}

func Update(v interface{}) error {
	return DB.Update(v)
}

func Delete(v interface{}) error {
	return DB.Delete(v)
}

func ForceDelete(v interface{}) error {
	return DB.ForceDelete(v)
}
