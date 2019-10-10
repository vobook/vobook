package factories

import (
	fake "github.com/brianvoe/gofakeit"

	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/enum/contact_property"
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
		m.Type = contactproperty.RandomType()
	}
	if m.Name == "" {
		m.Name = fake.Name()
	}
	if m.Value == "" {
		m.Name = fake.Company()
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
