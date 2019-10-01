package factories

import (
	fake "github.com/brianvoe/gofakeit"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/utils"
)

// MakeUser makes instance of user model
func MakeUser() (m models.User, err error) {
	email, err := utils.UniqueEmail("users")
	if err != nil {
		return
	}

	m = models.User{
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Email:     email,
		Password:  utils.RandomString(10),
	}

	return
}

// CreateUser makes instance of user model and inserts into DB
func CreateUser() (m models.User, err error) {
	m, err = MakeUser()
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
