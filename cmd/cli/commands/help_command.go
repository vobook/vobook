package commands

import (
	"fmt"
)

func init() {
	add("h", "help", command{
		handler: helpCommand,
		help:    `Display all available commands. Type h [command] to show help of given command`,
	})
}

func helpCommand(args ...string) (err error) {
	if len(args) > 0 {
		name := args[0]
		cmd, ok := repo[name]
		if ok {
			printCommandInfo(name, cmd.help)
			return
		} else {
			fmt.Printf("Command %s not found, try one of\n:", name)
		}
	}

	for name, cmd := range repo {
		printCommandInfo(name, cmd.help)
	}
	return
}

func printCommandInfo(name, text string) {
	fmt.Printf("%s - %s\n", name, text)
}
