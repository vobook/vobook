package emailverification

import (
	"vobook/cmd/server/errors"
	"vobook/database"
	"vobook/database/models"
	"vobook/utils"

	"github.com/go-pg/pg/v9"
)

func Create(userID, email string) (elem *models.EmailVerification, err error) {
	elem = &models.EmailVerification{
		UserID: userID,
		Email:  email,
	}

	token, err := utils.UniqueToken("email_verifications")
	if err != nil {
		return
	}
	elem.Token = token

	_, err = database.ORM().
		Model(elem).
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

func FindByToken(token string) (m models.EmailVerification, err error) {
	err = database.ORM().
		Model(&m).
		Where("token = ?", token).
		First()

	if err == pg.ErrNoRows {
		err = errors.EmailVerificationNotExists
		return
	}

	return
}

func DeleteByToken(id string) (err error) {
	_, err = database.ORM().
		Model(&models.EmailVerification{}).
		Where("token = ?", id).
		Delete()

	return
}
