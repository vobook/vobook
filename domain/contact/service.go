package contact

import (
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	contactproperty "github.com/vovainside/vobook/domain/contact_property"
)

func Create(m *models.Contact) (err error) {
	_, err = database.ORM().
		Model(m).
		Insert()
	if err != nil {
		return
	}

	for i := range m.Properties {
		m.Properties[i].ContactID = m.ID
	}

	err = contactproperty.CreateMany(&m.Properties)
	return
}

func Find(id string) (m models.Contact, err error) {
	err = database.ORM().
		Model(&m).
		Where("id = ?", id).
		First()

	return
}

func Update(m *models.Contact) (err error) {
	err = database.ORM().
		Update(m)

	return
}
