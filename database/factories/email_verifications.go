package factories

import (
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/utils"
)

func MakeEmailVerification() (m models.EmailVerification, err error) {
	userEl, err := CreateUser()
	if err != nil {
		return
	}

	email, err := utils.UniqueEmail("users")
	if err != nil {
		return
	}

	token, err := utils.UniqueToken("email_verifications")
	if err != nil {
		return
	}

	m = models.EmailVerification{
		Email:  email,
		UserID: userEl.ID,
		Token:  token,
	}

	return
}

func CreateEmailVerification() (m models.EmailVerification, err error) {
	m, err = MakeEmailVerification()
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
