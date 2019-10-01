package commands

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/go-pg/pg"
)

type command struct {
	help      string
	handler   func(...string) error
	countTime bool
}

var (
	repo = map[string]command{}
	DB   *pg.DB
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

func Run(name string, args ...string) error {
	cmd, ok := repo[name]
	if ok {
		startedAt := time.Now()
		err := cmd.handler(args...)
		if err != nil {
			return err
		}
		if cmd.countTime {
			duration := time.Since(startedAt)
			fmt.Printf("Done in %s\n", duration.Round(time.Millisecond).String())
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Unknown command %s, type help to see all available commands\n", name))
}
