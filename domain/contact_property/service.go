package contactproperty

import (
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
)

func CreateMany(elems []*models.ContactProperty) (err error) {
	if len(elems) == 0 {
		return
	}
	_, err = database.ORM().
		Model(elems).
		Insert()

	return
}
