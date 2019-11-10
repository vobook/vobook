package factories

import (
	"time"

	fake "github.com/brianvoe/gofakeit"

	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
)

func MakeContact(mOpt ...models.Contact) (m models.Contact, err error) {
	if len(mOpt) == 1 {
		m = mOpt[0]
	}

	if m.UserID == "" {
		var userEl models.User
		userEl, err = CreateUser()
		if err != nil {
			return
		}
		m.UserID = userEl.ID
	}
	if m.FirstName == "" {
		m.FirstName = fake.FirstName()
	}
	if m.LastName == "" {
		m.LastName = fake.LastName()
	}
	if m.MiddleName == "" {
		m.MiddleName = fake.LastName()
	}
	if m.Birthday == nil {
		birthday := fake.DateRange(time.Now().AddDate(-100, 0, 0), time.Now())
		m.Birthday = &birthday
	}

	return
}

func CreateContact(mOpt ...models.Contact) (m models.Contact, err error) {
	m, err = MakeContact(mOpt...)
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
