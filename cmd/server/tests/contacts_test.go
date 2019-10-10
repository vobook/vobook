package tests

import (
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit"

	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/database/factories"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/enum/contact_property"
	. "github.com/vovainside/vobook/tests/apitest"
	"github.com/vovainside/vobook/tests/assert"
	"github.com/vovainside/vobook/utils"
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

	var resp models.Contact
	TestCreate(t, "contacts", req, &resp)

	assert.Equals(t, AuthUser.ID, resp.UserID)
	assert.Equals(t, req.FirstName, resp.FirstName)
	assert.Equals(t, req.LastName, resp.LastName)
	assert.Equals(t, req.Birthday.Format(Conf.DateFormat), resp.Birthday.Format(Conf.DateFormat))
	assert.Equals(t, len(req.Properties), len(resp.Properties))

	for i, v := range req.Properties {
		assert.Equals(t, resp.ID, resp.Properties[i].ContactID)
		assert.Equals(t, v.Name, resp.Properties[i].Name)
		assert.Equals(t, v.Value, resp.Properties[i].Value)
		assert.Equals(t, v.Type, resp.Properties[i].Type)
	}

	assert.DatabaseHas(t, "contacts", utils.M{
		"id":         resp.ID,
		"user_id":    AuthUser.ID,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	})

	for _, v := range req.Properties {
		assert.DatabaseHas(t, "contact_properties", utils.M{
			"contact_id": resp.ID,
			"name":       v.Name,
			"value":      v.Value,
			"type":       v.Type,
		})
	}
}

func TestUpdateContact(t *testing.T) {
	elem, err := factories.CreateContact()
	assert.NotError(t, err)

	name := fake.Name()
	firstName := fake.FirstName()
	lastName := fake.LastName()
	middleName := fake.LastName()
	bday := fake.DateRange(time.Now().AddDate(-100, 0, 0), time.Now())
	req := requests.UpdateContact{
		Name:       &name,
		FirstName:  &firstName,
		LastName:   &lastName,
		MiddleName: &middleName,
		Birthday:   &bday,
	}

	var resp responses.Success
	TestUpdate(t, "contacts/"+elem.ID, req, &resp)

	assert.DatabaseHas(t, "contacts", utils.M{
		"id":          elem.ID,
		"name":        name,
		"first_name":  firstName,
		"last_name":   lastName,
		"middle_name": middleName,
	})
}
