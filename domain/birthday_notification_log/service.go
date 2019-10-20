package birthdaynotificationlog

import (
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
)

func Create(contactID string) (err error) {
	elem := models.BirthdayNotificationLog{
		ContactID: contactID,
	}

	_, err = database.ORM().
		Model(&elem).
		Insert()

	return
}
