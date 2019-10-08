package user

import (
	"time"

	"github.com/go-pg/pg"

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
		Where("email = ?", email).
		First()
	if err == pg.ErrNoRows {
		err = errors.UserByEmailNotFound
	}

	return
}

func FindByID(id string) (elem models.User, err error) {
	err = database.ORM().
		Model(&elem).
		Where("id = ?", id).
		First()

	return
}

func UpdatePassword(id, password string) (err error) {
	_, err = database.ORM().
		Model(&models.User{}).
		Where("id = ?", id).
		Set("password = ?", password).
		Update()

	return
}

func Delete(id string) (err error) {
	_, err = database.ORM().
		Model(&models.User{}).
		Where("id = ?", id).
		Set("deleted_at = ?", time.Now()).
		Update()

	return
}

func Restore(id string) (err error) {
	_, err = database.ORM().
		Model(&models.User{}).
		Where("id = ?", id).
		Set("deleted_at = null").
		Update()

	return
}

func Update(m *models.User) (err error) {
	err = database.ORM().
		Update(m)

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
			err = errors.EmailChangeEmailInUser
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
