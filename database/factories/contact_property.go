package factories

import (
	fake "github.com/brianvoe/gofakeit"

	"vobook/database"
	"vobook/database/models"
	contactpropertytype "vobook/enum/contact_property_type"
)

func MakeContactProperty(mOpt ...models.ContactProperty) (m models.ContactProperty, err error) {
	if len(mOpt) == 1 {
		m = mOpt[0]
	}

	if m.ContactID == "" {
		var contact models.Contact
		contact, err = CreateContact()
		if err != nil {
			return
		}
		m.ContactID = contact.ID
	}

	if m.Type == 0 {
		m.Type = contactpropertytype.Random()
	}
	if m.Name == "" {
		m.Name = fake.Name()
	}
	if m.Value == "" {
		m.Value = fake.Company()
	}

	return
}

func CreateContactProperty(mOpt ...models.ContactProperty) (m models.ContactProperty, err error) {
	m, err = MakeContactProperty(mOpt...)
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
