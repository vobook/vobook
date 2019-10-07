package tests

import (
	"testing"
	"time"

	"github.com/vovainside/vobook/cmd/server/responses"

	fake "github.com/brianvoe/gofakeit"

	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/enum/contact_property"
	. "github.com/vovainside/vobook/tests/apitest"
)

func TestCreateContact(t *testing.T) {
	req := requests.CreateContact{
		Name:      fake.Name(),
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Birthday:  fake.DateRange(time.Now().AddDate(-100, 0, 0), time.Now()),
		Properties: []requests.CreateContactProperty{
			{
				Type:  contactproperty.TypeEmail,
				Value: fake.Email(),
			},
			{
				Type:  contactproperty.TypePhone,
				Value: fake.PhoneFormatted(),
			},
			{
				Type:  contactproperty.TypePhone,
				Value: fake.PhoneFormatted(),
				Name:  "Only SMS",
			},
			{
				Type:  contactproperty.TypeLink,
				Value: fake.URL(),
				Name:  "limit for me",
			},
			{
				Type:  contactproperty.TypeOther,
				Value: fake.Sentence(3),
				Name:  "hello",
			},
		},
	}

	var resp responses.Success
	TestCreate(t, "contacts", req, &resp)

	// TODO assert database has
}
