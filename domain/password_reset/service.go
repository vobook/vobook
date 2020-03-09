package passwordreset

import (
	"vobook/cmd/server/errors"
	"vobook/database"
	"vobook/database/models"
	"vobook/utils"

	"github.com/go-pg/pg/v9"
)

func Create(elem *models.PasswordReset) (err error) {
	token, err := utils.UniqueToken("password_resets")
	if err != nil {
		return
	}
	elem.Token = token

	_, err = database.ORM().
		Model(elem).
		Insert()

	return
}

func Find(token string) (elem models.PasswordReset, err error) {
	err = database.ORM().
		Model(&elem).
		Where("token = ?", token).
		First()
	if err == pg.ErrNoRows {
		err = errors.PasswordResetTokenNotFound
	}

	return
}
