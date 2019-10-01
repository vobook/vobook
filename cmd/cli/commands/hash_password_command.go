package commands

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	add("hash-password", "hp", command{
		handler: hashPassword,
		help: `Makes a hash of given password.
	Usage: hash-password [password]`,
	})
}

func hashPassword(args ...string) (err error) {
	if len(args) == 0 {
		fmt.Println("Password must be provided")
		return
	}

	password := args[0]

	// TODO make random password if not provided
	// TODO check for short password and warn
	// TODO check for simple password and warn

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(bytes))
	return
}
