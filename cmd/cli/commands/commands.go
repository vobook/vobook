package commands

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/vovainside/vobook/database"
)

type command struct {
	help      string
	handler   func(...string) error
	countTime bool
}

var (
	repo = map[string]command{}
	DB   orm.DB
)

func add(params ...interface{}) {
	var names = []string{}
	var cmd command

	for _, p := range params {
		t := reflect.TypeOf(p).String()

		if t == "string" {
			names = append(names, p.(string))
			continue
		}

		if t == "commands.command" {
			cmd = p.(command)
			continue
		}

		fmt.Printf("Wrong param type given for AddCommand, must be eiher (string) or (main.Command), but (%s) given\n", t)
	}

	for _, name := range names {
		_, ok := repo[name]
		if ok {
			log.Fatalf("Command [%s] already registered", name)
		}
		repo[name] = cmd
	}
}

func Run(name string, args ...string) (err error) {
	cmd, ok := repo[name]
	if ok {
		db := database.Conn()
		var tx *pg.Tx
		tx, err = db.Conn().Begin()
		if err != nil {
			return
		}
		database.SetDB(tx)
		DB = database.ORM()

		startedAt := time.Now()
		err = cmd.handler(args...)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				println("tx rollback err:", err2.Error())
			}
			return err
		}

		if cmd.countTime {
			duration := time.Since(startedAt)
			fmt.Printf("Done in %s\n", duration.Round(time.Millisecond).String())
		}

		err2 := tx.Commit()
		if err2 != nil {
			println("tx commit err:", err2.Error())
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Unknown command %s, type help to see all available commands\n", name))
}
