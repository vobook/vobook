package emailverification

import (
	"github.com/go-pg/pg"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
)

func Create(userID, email string) (err error) {
	elem := models.EmailVerification{
		UserID: userID,
		Email:  email,
	}
	_, err = database.ORM().
		Model(&elem).
		Insert()
	return
}

func FindByID(id string) (m models.EmailVerification, err error) {
	err = database.ORM().
		Model(&m).
		Where("id = ?", id).
		First()

	if err == pg.ErrNoRows {
		err = errors.EmailVerificationNotExists
		return
	}

	return
}

func Delete(id string) (err error) {
	_, err = database.ORM().
		Model(&models.EmailVerification{}).
		Where("id = ?", id).
		Delete()

	return
}
