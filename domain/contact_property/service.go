package contactproperty

import (
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
)

func CreateMany(elems *[]models.ContactProperty) (err error) {
	if len(*elems) == 0 {
		return
	}
	_, err = database.ORM().
		Model(elems).
		Insert()

	return
}

func Find(id string) (elem models.ContactProperty, err error) {
	err = database.ORM().
		Model(&elem).
		Where("id = ?", id).
		First()

	return
}

func Update(elem *models.ContactProperty) (err error) {
	err = database.ORM().
		Update(elem)

	return
}
