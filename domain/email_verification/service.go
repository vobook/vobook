package emailverification

import (
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
