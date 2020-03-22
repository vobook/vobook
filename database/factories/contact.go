package factories

import (
	"time"

	fake "github.com/brianvoe/gofakeit"

	"vobook/database"
	"vobook/database/models"
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

	// Birthday might be:
	// unknown, only year, only month and day, full date
	if m.DOBYear == 0 && m.DOBMonth == 0 && m.DOBDay == 0 && fake.Bool() {
		if fake.Bool() {
			m.DOBYear = fake.Year()
		}
		if fake.Bool() {
			m.DOBMonth = time.Month(fake.Number(1, 12))
			m.DOBDay = fake.Day()
		}
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
