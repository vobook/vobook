package user

import (
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
)

func Create(elem *models.User) (err error) {
	_, err = database.ORM().
		Model(elem).
		Insert()

	return
}

func FindByEmail(email string) (elem models.User, err error) {
	err = database.ORM().
		Model(&elem).
		Where("email=?", email).
		First()

	return
}
