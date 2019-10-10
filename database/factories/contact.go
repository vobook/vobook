package factories

import (
	"time"

	fake "github.com/brianvoe/gofakeit"

	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
)

func MakeContact() (m models.Contact, err error) {
	userEl, err := CreateUser()
	if err != nil {
		return
	}

	m = models.Contact{
		UserID:     userEl.ID,
		FirstName:  fake.FirstName(),
		LastName:   fake.LastName(),
		MiddleName: fake.LastName(),
		Birthday:   fake.DateRange(time.Now().AddDate(-100, 0, 0), time.Now()),
	}

	return
}

func CreateContact() (m models.Contact, err error) {
	m, err = MakeContact()
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
