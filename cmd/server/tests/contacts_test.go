package tests

import (
	"fmt"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit"
	"github.com/davecgh/go-spew/spew"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/factories"
	"github.com/vovainside/vobook/database/models"
	contactpropertytype "github.com/vovainside/vobook/enum/contact_property_type"
	. "github.com/vovainside/vobook/tests/apitest"
	"github.com/vovainside/vobook/tests/assert"
	"github.com/vovainside/vobook/utils"
)

func TestCreateContact(t *testing.T) {
	req := requests.CreateContact{
		Name:      fake.Name(),
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Birthday:  fake.DateRange(time.Now().AddDate(-100, 0, 0), time.Now()).Format(database.DateFormat),
		Properties: []requests.CreateContactProperty{
			{
				Type:  contactpropertytype.Email,
				Value: fake.Email(),
			},
			{
				Type:  contactpropertytype.Phone,
				Value: fake.PhoneFormatted(),
			},
			{
				Type:  contactpropertytype.Phone,
				Value: fake.PhoneFormatted(),
				Name:  "Only SMS",
			},
			{
				Type:  contactpropertytype.Link,
				Value: fake.URL(),
				Name:  "limit for me",
			},
			{
				Type:  contactpropertytype.Other,
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
	assert.Equals(t, req.Birthday, resp.Birthday.Format(Conf.DateFormat))
	assert.Equals(t, len(req.Properties), len(resp.Props))

	for i, v := range req.Properties {
		assert.Equals(t, resp.ID, resp.Props[i].ContactID)
		assert.Equals(t, v.Name, resp.Props[i].Name)
		assert.Equals(t, v.Value, resp.Props[i].Value)
		assert.Equals(t, v.Type, resp.Props[i].Type)
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
	Login(t)
	elem, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	name := fake.Name()
	firstName := fake.FirstName()
	lastName := fake.LastName()
	middleName := fake.LastName()
	bday := fake.DateRange(time.Now().AddDate(-100, 0, 0), time.Now()).Format(config.Get().DateFormat)
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

func TestTrashContacts(t *testing.T) {
	Login(t)
	elem1, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)
	elem2, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)
	elem3, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	req := requests.IDs{
		elem1.ID,
		elem2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "trash-contacts", req, &resp)

	assert.DatabaseHasDeleted(t, "contacts", elem1.ID, elem2.ID)
	assert.DatabaseHas(t, "contacts", utils.M{
		"id":         elem3.ID,
		"deleted_at": nil,
	})
}

func TestRestoreContacts(t *testing.T) {
	Login(t)
	deletedAt := time.Now()
	elem1, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID, DeletedAt: &deletedAt})
	assert.NotError(t, err)
	elem2, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID, DeletedAt: &deletedAt})
	assert.NotError(t, err)
	elem3, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID, DeletedAt: &deletedAt})
	assert.NotError(t, err)

	req := requests.IDs{
		elem1.ID,
		elem2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "restore-contacts", req, &resp)

	assert.DatabaseHasDeleted(t, "contacts", elem3.ID)
	assert.DatabaseHas(t, "contacts", utils.M{
		"id":         elem1.ID,
		"deleted_at": nil,
	})
	assert.DatabaseHas(t, "contacts", utils.M{
		"id":         elem2.ID,
		"deleted_at": nil,
	})
}

func TestDeleteContacts(t *testing.T) {
	Login(t)
	elem1, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)
	elem2, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)
	elem3, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	req := requests.IDs{
		elem1.ID,
		elem2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "delete-contacts", req, &resp)

	assert.DatabaseHas(t, "contacts", utils.M{
		"id": elem3.ID,
	})
	assert.DatabaseMissing(t, "contacts", utils.M{
		"id": elem1.ID,
	})
	assert.DatabaseMissing(t, "contacts", utils.M{
		"id": elem2.ID,
	})
}

func TestGetContact(t *testing.T) {
	Login(t)
	elem, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	props := make([]models.ContactProperty, 3)
	for i := range props {
		prop, err := factories.CreateContactProperty(models.ContactProperty{
			ContactID: elem.ID,
			Name:      fmt.Sprintf("Prop %d", i+1),
			Order:     i + 1,
		})
		assert.NotError(t, err)
		props[i] = prop
	}

	var resp models.Contact
	Fetch(t, "contacts/"+elem.ID, &resp)

	assert.Equals(t, elem.ID, resp.ID)
	assert.Equals(t, elem.FirstName, resp.FirstName)
	assert.Equals(t, elem.LastName, resp.LastName)
	assert.Equals(t, 3, len(resp.Props))

	spew.Dump(resp)

	for i, v := range props {
		assert.Equals(t, v.ID, resp.Props[i].ID)
		assert.Equals(t, v.Type, resp.Props[i].Type)
		assert.Equals(t, v.Name, resp.Props[i].Name)
		assert.Equals(t, v.Value, resp.Props[i].Value)
		assert.Equals(t, v.Order, resp.Props[i].Order)
	}
}

func TestSearchContact(t *testing.T) {
	ReLogin(t)
	elem, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	props := make([]models.ContactProperty, 3)
	for i := range props {
		prop, err := factories.CreateContactProperty(models.ContactProperty{
			ContactID: elem.ID,
			Name:      fmt.Sprintf("Prop %d", i+1),
			Order:     i + 1,
		})
		assert.NotError(t, err)
		props[i] = prop
	}

	_, err = factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	deletedAt := time.Now()
	elem3, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID, DeletedAt: &deletedAt})
	assert.NotError(t, err)

	// should get 2 contacts
	var req requests.SearchContact
	var resp responses.SearchContact
	TestSearch(t, "contacts", req, &resp)
	assert.Equals(t, 2, len(resp.Data))

	// should get 1 deleted contact
	req = requests.SearchContact{Trashed: true}
	TestSearch(t, "contacts", req, &resp)
	assert.Equals(t, 1, len(resp.Data))
	assert.Equals(t, elem3.ID, resp.Data[0].ID)

	// should get 1 filtered contact
	req = requests.SearchContact{Query: props[0].Value}
	TestSearch(t, "contacts", req, &resp)
	assert.Equals(t, 1, len(resp.Data))
	assert.Equals(t, elem.ID, resp.Data[0].ID)
}
