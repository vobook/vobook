package user

import (
	"github.com/vovainside/vobook/cmd/server/errors"
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

func FindByID(id string) (elem models.User, err error) {
	err = database.ORM().
		Model(&elem).
		Where("id=?", id).
		First()

	return
}

func EmailVerified(id, email string) (err error) {
	if email != "" {
		var count int
		count, err = database.ORM().Model(&models.User{}).
			Where("email = ?", email).
			Where("id != ?", id).
			Count()
		if err != nil {
			return
		}
		if count > 0 {
			err = errors.EmailChangeEmailAlreadyExists
			return
		}
	}

	q := database.ORM().
		Model(&models.User{}).
		Set("email_verified = true").
		Where("id = ?", id)

	if email != "" {
		q.Set("email = ?", email)
	}

	_, err = q.Update()
	return
}
