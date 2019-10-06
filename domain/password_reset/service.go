package passwordreset

import (
	"github.com/go-pg/pg"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/utils"
)

func Create(elem *models.PasswordReset) (err error) {
	elem.Token = utils.RandomHash()

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
