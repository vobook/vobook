package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vovainside/vobook/cmd/cli/commands"
)

const usageText = `
VOBOOK CLI

This program runs various commands to make dev's life simpler
Enter h or help to list all available commands
Enter h [command] to show description of given command`

func main() {
	if len(os.Args) > 1 {
		args := os.Args[2:]
		err := commands.Run(os.Args[1], args...)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(usageText)

	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for {
		fmt.Print("> ")
		scanner.Scan()
		input = scanner.Text()
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")
		name := args[0]
		if len(args) > 1 {
			args = args[1:]
		} else {
			args = []string{}
		}

		err := commands.Run(name, args...)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
