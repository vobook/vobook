package birthdaynotificationlog

import (
	"vobook/database"
	"vobook/database/models"
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
