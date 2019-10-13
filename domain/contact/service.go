package contact

import (
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/domain/contact_property"
)

func Create(m *models.Contact) (err error) {
	_, err = database.ORM().
		Model(m).
		Insert()
	if err != nil {
		return
	}

	for i := range m.Props {
		m.Props[i].ContactID = m.ID
	}

	err = contactproperty.CreateMany(&m.Props)
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

func Props(id string) (elems []models.ContactProperty, err error) {
	err = database.ORM().
		Model(&elems).
		Where("contact_id = ?", id).
		Order("order ASC").
		Select()

	return
}
