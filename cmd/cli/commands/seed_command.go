package commands

import (
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/utils"
)

func init() {
	add("seed", "sd", command{
		handler: seed,
		help:    `Seeds database with test data. Should be run on empty database`,
	})
}

func seed(args ...string) (err error) {
	passwordHash, err := utils.HashPassword("test")
	if err != nil {
		return
	}

	// Users
	users := []models.User{
		{
			FirstName:     "John",
			LastName:      "Snow",
			Email:         "test@test.me",
			EmailVerified: true,
		},
		{
			FirstName:     "Vladimir",
			LastName:      "Ognev",
			Email:         "vo@proj.run",
			EmailVerified: true,
		},
	}

	for i := range users {
		users[i].Password = passwordHash
	}

	err = DB.Insert(&users)
	if err != nil {
		return
	}
	println(len(users), "users created")

	return
}
